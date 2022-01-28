package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/samuelventura/go-modbus"
	"github.com/samuelventura/go-serial"
)

func powerOn(secondsOn int, master modbus.Master) {
	for i := 0; i < secondsOn; i++ {
		err := master.WriteDo(1, 4096, true)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(".")
		time.Sleep(time.Second)
	}
	fmt.Println("")
}

func powerOff(master modbus.Master) {
	err := master.WriteDo(1, 4096, false)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(6 * time.Second)
}

//(cd sample; go run .)
func main() {
	scheme := ""
	ip := ""
	start := ""
	end := ""
	if len(os.Args) == 5 {
		start = os.Args[1]
		end = os.Args[2]
		scheme = os.Args[3]
		ip = os.Args[4]
	} else {
		log.Fatal("usage: lvfix starpowerOfft end scheme ip")
	}
	startInt, err := strconv.Atoi(start)
	if err != nil {
		log.Fatal(err)
	}
	endInt, err := strconv.Atoi(end)
	if err != nil {
		log.Fatal(err)
	}

	log.SetFlags(log.Lmicroseconds)
	mode := &serial.Mode{}
	mode.BaudRate = 9600
	mode.DataBits = 8
	mode.Parity = serial.NoParity
	mode.StopBits = serial.OneStopBit
	trans, err := serial.NewSerialTransport("/dev/ttyUSB0", mode)
	if err != nil {
		log.Fatal(err)
		println("Missing device")
	}
	defer trans.Close()
	trans.DiscardOn()
	modbus.EnableTrace(false)
	master := modbus.NewRtuMaster(trans, 400)

	counter := startInt

	powerOff(master)

	for {
		secondsOn := counter
		fmt.Printf("%v", secondsOn)
		powerOn(secondsOn, master)
		printWebSocket(scheme, ip, secondsOn)
		powerOff(master)
		counter++
		if counter > endInt {
			counter = startInt
		}
	}
}
