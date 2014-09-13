package main

import (
	"flag"
	"fmt"
	"net"
	"bytes"
	"time"
	"log"
	"encoding/binary"
)

type Message struct {
	Sender [20]byte
	DateTime int64
	Message [100]byte
}

func (m *Message) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, m)
	return buf.Bytes(), err
}

func (m *Message) UnmarshalBinary(data []byte) error {
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.BigEndian, m)
	return err
}

func broadcast(name string, message string) {
	packet := Message{DateTime: time.Now().Unix()}
	copy(packet.Sender[:], name)
	copy(packet.Message[:], message)

	conn, err := net.Dial("udp", "255.255.255.255:9666")
	if err != nil {
		log.Fatal(err)
	}
	
	packetBytes, _ := packet.MarshalBinary()
	conn.Write(packetBytes)
}

func receive() {
	addr, _ := net.ResolveUDPAddr("udp", ":9666")
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal("Could not connect!")
	}

	for {
		buf := make([]byte, 128)
		conn.Read(buf)
		var message Message
		message.UnmarshalBinary(buf)
		fmt.Printf("[%d]%s: %s\n", message.DateTime, string(message.Sender[:]), string(message.Message[:]))
	}
}

func main () {
	name := flag.String("name", "anon", "Your username")
	message := flag.String("message", "Just connected", "Your message")

	flag.Parse()
	fmt.Print(flag.Args())
	if len(flag.Args()) == 0 {
		fmt.Println("Sending...")
		broadcast(*name, *message)
	} else {
		fmt.Println("Receiving...")
		receive()
	}
}

