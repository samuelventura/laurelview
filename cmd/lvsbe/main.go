package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	go func() {
		defer os.Exit(0)
		ioutil.ReadAll(os.Stdin)
	}() //exit on stdin close
	endpoint := endpoint()
	gin.SetMode(gin.ReleaseMode) //remove debug warning
	router := gin.New()          //remove default logger
	router.Use(gin.Recovery())   //catches panics
	router.Use(static)
	router.GET("/discovery/:tos", func(c *gin.Context) {
		toss := c.Param("tos")
		tos, err := strconv.ParseInt(toss, 10, 32)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		if tos <= 0 || tos >= 5 {
			c.String(http.StatusBadRequest, "invalid tos (0, 5)")
			return
		}
		list, err := discover(int(tos))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, list)
	})
	router.GET("/blink/:ips", func(c *gin.Context) {
		ips := c.Param("ips")
		err := blink(ips)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, "ok")
	})
	//FIXME
	listen, err := net.Listen("tcp", endpoint)
	if err != nil {
		log.Fatal(err)
	}
	server := &http.Server{
		Addr:    endpoint,
		Handler: router,
	}
	err = server.Serve(listen)
	if err != nil {
		log.Println(err)
	}
}

func endpoint() string {
	ep := os.Getenv("LV_SBE_ENDPOINT")
	if len(strings.TrimSpace(ep)) > 0 {
		return ep
	}
	return ":0"
}

//go:embed build/*
var build embed.FS

func static(c *gin.Context) {
	path := "build" + c.Request.URL.Path
	if path == "build/" {
		path = "build/index.html"
	}
	data, err := build.ReadFile(path)
	if err != nil {
		c.Next()
	} else {
		ext := filepath.Ext(path)
		ct := mime.TypeByExtension(ext)
		c.Data(http.StatusOK, ct, data)
	}
}

type IdRequestDso struct {
	Name   string `json:"name"`
	Action string `json:"action"`
}

type IdResponseDso struct {
	Name   string         `json:"name"`
	Action string         `json:"action"`
	Data   IdResponseData `json:"data"`
}

type IdResponseData struct {
	Hostname string `json:"hostname"`
	Ifname   string `json:"ifname"`
	MacAddr  string `json:"macaddr"`
	Name     string `json:"name"`
	Version  string `json:"version"`
	IpFrom   string `json:"ipfrom"`
	IpAddr   string `json:"ipaddr"`
}

func discover(tos int) ([]*IdResponseDso, error) {
	socket, err := net.ListenUDP("udp4", &net.UDPAddr{})
	if err != nil {
		return nil, err
	}
	defer socket.Close()
	log.Println("LocalAddr", socket.LocalAddr())
	idb, err := json.Marshal(&IdRequestDso{
		Name: "lvbox", Action: "id"})
	if err != nil {
		return nil, err
	}
	log.Println(">", string(idb))
	idn, err := socket.WriteToUDP(idb, &net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255),
		Port: 31680,
	})
	if err != nil || idn != len(idb) {
		return nil, err
	}
	list := []*IdResponseDso{}
	inbuf := make([]byte, 2048)
	tosd := time.Duration(tos)
	socket.SetDeadline(time.Now().Add(tosd * time.Second))
	for {
		inn, addr, err := socket.ReadFromUDP(inbuf)
		nerr, ok := err.(net.Error)
		if ok && nerr.Timeout() {
			return list, nil
		}
		if err != nil {
			return nil, err
		}
		log.Println("<", addr, string(inbuf[:inn]))
		response := &IdResponseDso{}
		err = json.Unmarshal(inbuf[:inn], response)
		if err != nil {
			log.Println(err)
		} else {
			response.Data.IpFrom = addr.IP.String()
			log.Println(response)
			list = append(list, response)
		}
	}
}

func blink(ips string) error {
	ip := net.ParseIP(ips)
	if ip == nil {
		return fmt.Errorf("invalid ip")
	}
	socket, err := net.ListenUDP("udp4", &net.UDPAddr{})
	if err != nil {
		return err
	}
	defer socket.Close()
	log.Println("LocalAddr", socket.LocalAddr())
	idb, err := json.Marshal(&IdRequestDso{
		Name: "lvbox", Action: "blink"})
	if err != nil {
		return err
	}
	log.Println(">", string(idb))
	idn, err := socket.WriteToUDP(idb, &net.UDPAddr{
		IP:   ip,
		Port: 31680,
	})
	if err != nil || idn != len(idb) {
		return err
	}
	return nil
}
