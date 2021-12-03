package main

import (
	"encoding/json"
	"log"
	"net"
	"time"
)

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
	IpAddr   string `json:"ipaddr"`
}

func main() {
	socket, err := net.ListenUDP("udp4", &net.UDPAddr{})
	if err != nil {
		log.Fatal(err)
	}
	defer socket.Close()
	log.Println("LocalAddr", socket.LocalAddr())
	idb, err := json.Marshal(&IdRequestDso{
		Name: "nerves", Action: "id"})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(">", string(idb))
	idn, err := socket.WriteToUDP(idb, &net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255),
		Port: 31680,
	})
	if err != nil || idn != len(idb) {
		log.Fatal(err)
	}
	inbuf := make([]byte, 2048)
	socket.SetDeadline(time.Now().Add(1 * time.Second))
	for {
		inn, addr, err := socket.ReadFromUDP(inbuf)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("<", addr, string(inbuf[:inn]))
		response := &IdResponseDso{}
		err = json.Unmarshal(inbuf[:inn], response)
		if err != nil {
			log.Println(err)
		} else {
			response.Data.IpAddr = addr.IP.String()
			log.Println(response)
		}
	}
}
