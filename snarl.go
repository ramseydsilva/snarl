package snarl

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

func Broadcast(name string, message string, domain string, port int) {
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

func Receive(port int) chan Message {
	addr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal("Could not connect")
	}

	messageChannel := make(chan Message)
	go func() {
		for {
			var message Message
			buf := make([]byte, 128)
			conn.Read(buf)
			message.UnmarshalBinary(buf)
			messageChannel <- message
		}
	}()
	return messageChannel
}

func main() {
	port := flag.Int("port", 9666, "Private channel")
	domain := flag.String("domain", "255.255.255.255", "Domain")
	name := flag.String("name", "anon", "Your username")
	message := flag.String("message", "Just connected", "Your message")

	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Println("Sending...")
		Broadcast(*name, *message, *domain, *port)
	} else {
		fmt.Println("Receiving...")
		for m := range Receive(*port) {
			fmt.Printf("[%d] %s: %s\n", m.DateTime, string(m.Sender[:]), string(m.Message[:]))
		}
	}
}
