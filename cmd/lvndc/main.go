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
}

func main() {
	listen, err := net.ListenUDP("udp4", &net.UDPAddr{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("LocalAddr", listen.LocalAddr())
	id := &IdRequestDso{Name: "nerves", Action: "id"}
	idb, err := json.Marshal(id)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(">", string(idb))
	idn, err := listen.WriteToUDP(idb, &net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255),
		Port: 31680,
	})
	if err != nil || idn != len(idb) {
		log.Fatal(err)
	}
	input := make([]byte, 2048)
	listen.SetDeadline(time.Now().Add(1 * time.Second))
	for {
		inn, _, err := listen.ReadFromUDP(input)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("<", string(input[:inn]))
		response := &IdResponseDso{}
		err = json.Unmarshal(input[:inn], response)
		if err != nil {
			log.Println(err)
		} else {
			log.Println(response)
		}
	}
}
