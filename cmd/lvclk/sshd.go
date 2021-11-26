package main

import (
	"embed"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/samuelventura/go-tools"
	"github.com/samuelventura/go-tree"
	"golang.org/x/crypto/ssh"
)

//go:embed keys/*
var keys embed.FS

func sshd(node tree.Node) {
	endpoint := node.GetValue("endpoint").(string)
	publicBytes, err := keys.ReadFile("keys/id_rsa.pub")
	if err != nil {
		log.Panicln(err)
	}
	privateBytes, err := keys.ReadFile("keys/id_rsa")
	if err != nil {
		log.Panicln(err)
	}
	pubkey, _, _, _, err := ssh.ParseAuthorizedKey(publicBytes)
	if err != nil {
		log.Panicln(err)
	}
	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Panicln(err)
	}
	pubtxt := strings.TrimSpace(string(ssh.MarshalAuthorizedKey(pubkey)))
	config := &ssh.ServerConfig{
		PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			inkey := strings.TrimSpace(string(ssh.MarshalAuthorizedKey(key)))
			if pubtxt == inkey {
				return &ssh.Permissions{Extensions: map[string]string{"key-id": "defkey"}}, nil
			}
			return nil, fmt.Errorf("key not found")
		},
	}
	config.AddHostKey(private)
	node.SetValue("config", config)
	listen, err := net.Listen("tcp", endpoint)
	if err != nil {
		log.Panicln(err)
	}
	node.AddCloser("listen", listen.Close)
	port := listen.Addr().(*net.TCPAddr).Port
	node.SetValue("port", port)
	node.AddProcess("listen", func() {
		id := NewId("ssh-" + listen.Addr().String())
		for {
			tcpConn, err := listen.Accept()
			if err != nil {
				log.Println(err)
				return
			}
			setupSshConnection(node, tcpConn, id)
		}
	})
}

func setupSshConnection(node tree.Node, tcpConn net.Conn, id Id) {
	defer node.IfRecoverCloser(tcpConn.Close)
	addr := tcpConn.RemoteAddr().String()
	cid := id.Next(addr)
	child := node.AddChild(cid)
	child.AddCloser("tcpConn", tcpConn.Close)
	child.AddProcess("tcpConn", func() {
		handleSshConnection(child, tcpConn)
	})
}

func handleSshConnection(node tree.Node, tcpConn net.Conn) {
	tools.KeepAlive(tcpConn, 5)
	single := node.GetValue("single").(*singleDso)
	proxy := node.GetValue("proxy").(string)
	config := node.GetValue("config").(*ssh.ServerConfig)
	sshConn, chans, reqs, err := ssh.NewServerConn(tcpConn, config)
	if err != nil {
		log.Println(err)
		return
	}
	node.AddCloser("sshConn", sshConn.Close)
	node.SetValue("ssh", sshConn)
	single.enter(node)
	defer single.exit(node)
	listen, err := net.Listen("tcp", proxy)
	if err != nil {
		log.Println(err)
		return
	}
	node.AddCloser("listen", listen.Close)
	log.Println(proxy, tcpConn.RemoteAddr())
	node.AddProcess("ssh chans reject", func() {
		for nch := range chans {
			nch.Reject(ssh.Prohibited, "unsupported")
		}
	})
	node.AddProcess("ssh reqs reply", func() {
		for req := range reqs {
			if req.WantReply {
				req.Reply(false, nil)
			}
		}
	})
	node.AddProcess("ssh ping handler", func() {
		for {
			dl := time.Now().Add(10 * time.Second)
			resp, _, err := sshConn.SendRequest("ping", true, nil)
			if time.Now().After(dl) || err != nil || !resp {
				log.Println(proxy, "ping timeout")
				return
			}
			timer := time.NewTimer(5 * time.Second)
			select {
			case <-timer.C:
				continue
			case <-node.Closed():
				timer.Stop()
				return
			}
		}
	})
	id := NewId("proxy-" + listen.Addr().String())
	for {
		proxyConn, err := listen.Accept()
		if err != nil {
			log.Println(proxy, err)
			break
		}
		setupProxyConnection(node, proxyConn, id)
	}
}

func setupProxyConnection(node tree.Node, proxyConn net.Conn, id Id) {
	defer node.IfRecoverCloser(proxyConn.Close)
	addr := proxyConn.RemoteAddr().String()
	cid := id.Next(addr)
	child := node.AddChild(cid)
	child.AddCloser("proxyConn", proxyConn.Close)
	child.AddProcess("proxyConn", func() {
		handleProxyConnection(child, proxyConn)
	})
}

func handleProxyConnection(node tree.Node, proxyConn net.Conn) {
	tools.KeepAlive(proxyConn, 5)
	proxy := node.GetValue("proxy").(string)
	sshConn := node.GetValue("ssh").(*ssh.ServerConn)
	sshChan, reqChan, err := sshConn.OpenChannel("forward", nil)
	if err != nil {
		log.Println(proxy, err)
		return
	}
	node.AddCloser("sshChan", sshChan.Close)
	node.AddProcess("DiscardRequests(reqChan)", func() {
		ssh.DiscardRequests(reqChan)
	})
	node.AddProcess("Copy(sshChan, proxyConn)", func() {
		_, err := io.Copy(sshChan, proxyConn)
		if err != nil {
			log.Println(proxy, err)
		}
	})
	node.AddProcess("Copy(proxyConn, sshChan)", func() {
		_, err := io.Copy(proxyConn, sshChan)
		if err != nil {
			log.Println(proxy, err)
		}
	})
	node.WaitClosed()
}
