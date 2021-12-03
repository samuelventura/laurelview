package main

import (
	"encoding/json"
	"log"
	"net"
	"time"
)

type IdDso struct {
	Name   string `json:"name"`
	Action string `json:"action"`
}

func main() {
	ra, err := net.ResolveUDPAddr("udp4", "255.255.255.255:31680")
	if err != nil {
		log.Fatal(err)
	}
	la, err := net.ResolveUDPAddr("udp4", "0.0.0.0:0")
	if err != nil {
		log.Fatal(err)
	}
	listen, err := net.ListenUDP("udp4", la)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("LocalAddr", listen.LocalAddr())
	id := &IdDso{Name: "nerves", Action: "id"}
	idb, err := json.Marshal(id)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(">", string(idb))
	idn, err := listen.WriteToUDP(idb, ra)
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
	}
}
