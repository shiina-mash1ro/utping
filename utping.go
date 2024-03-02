package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
)

func main() {
	var address string

	args := os.Args
	argc := len(os.Args)
	if argc < 2 {
		fmt.Println("usage: utping example.net:12345")
		os.Exit(0)
	}

	address = args[1]

	if address == "" {
		panic("address is required")
	}

	connId := uint16(rand.Int())

	sendHeader := header{
		Type:          4,
		Version:       1,
		ConnID:        connId,
		SeqNr:         1,
		AckNr:         0,
		WndSize:       0,
		Timestamp:     0,
		TimestampDiff: 0,
	}

	buff := make([]byte, 40)
	sendHeader.Marshal(buff)

	receiveHeader := header{}

	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	dstAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		panic(err)
	}
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	time1 := time.Now()
	conn.Write(buff)

	data := make([]byte, 1024)

	n, err := conn.Read(data)
	time2 := time.Now()
	receiveHeader.Unmarshal(data)

	fmt.Print(receiveHeader)
	fmt.Printf("read %s from <%s>\n", data[:n], conn.RemoteAddr())

	if receiveHeader.Type == 2 {
		fmt.Printf("latency: %d ms\n", time2.Sub(time1).Milliseconds())
		os.Exit(0)
	} else {
		fmt.Printf("fail\n")
		os.Exit(1)
	}
}
