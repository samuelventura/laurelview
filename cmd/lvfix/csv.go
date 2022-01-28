package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"
)

var path = "result.csv"

//var data = [][]string{{"Time Stamp", "Event"}}

func writeFile(seconds int, result string) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)

	if err != nil {
		log.Fatal(err)
		println("File error")
	}

	defer file.Close()

	var data [][]string
	timeStamp := time.Now()
	newData := result
	layout := "2006-01-02 15:04:05"
	data = append(data, []string{timeStamp.Format(layout), strconv.Itoa(seconds), newData})

	writer := csv.NewWriter(file)
	writer.WriteAll(data)

	if err != nil {
		log.Fatal(err)
	}
}
