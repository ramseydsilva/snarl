package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

type Message struct {
	Sender   [20]byte
	DateTime int64
	Message  [100]byte
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

func broadcast(name string, message string, domain string, port int) {
	packet := Message{DateTime: time.Now().Unix()}
	copy(packet.Sender[:], name)
	copy(packet.Message[:], message)

	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", domain, port))
	if err != nil {
		log.Fatal(err)
	}

	packetBytes, _ := packet.MarshalBinary()
	conn.Write(packetBytes)
}

func receive(port int) {
	addr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
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

func main() {
	port := flag.Int("port", 9666, "Private channel")
	domain := flag.String("domain", "255.255.255.255", "Domain")
	name := flag.String("name", "anon", "Your username")
	message := flag.String("message", "Just connected", "Your message")

	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Println("Sending...")
		broadcast(*name, *message, *domain, *port)
	} else {
		fmt.Println("Receiving...")
		receive(*port)
	}
}
