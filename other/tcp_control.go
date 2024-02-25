package main

import (
	"fmt"
	"log"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"golang.org/x/net/ipv4"
)

func main() {
	proto := layers.IPProtocolTCP

	c, err := net.ListenPacket(fmt.Sprintf("ip4:%d", proto), "127.0.0.1")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	iprawconn, err := ipv4.NewRawConn(c)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("iprawconn:%+v\n", iprawconn)

	tcp := layers.TCP{
		SrcPort: 43124,
		DstPort: 2222,
		SYN:     true,
		ACK:     false,
		Window:  65535,
		//Checksum:
		Seq: 231,
	}
	// payload := []byte{'a', 'b', 'c', '\n'}
	options := gopacket.SerializeOptions{
		// ComputeChecksums: true,
		FixLengths: true,
	}
	buffer := gopacket.NewSerializeBuffer()
	err = gopacket.SerializeLayers(buffer, options,
		&tcp,
		// gopacket.Payload(payload),
	)
	if err != nil {
		fmt.Printf("[-] Serialize error: %s\n", err.Error())
		return
	}
	b := buffer.Bytes()
	fmt.Println(b)

	h := &ipv4.Header{
		Version:  ipv4.Version,
		Len:      ipv4.HeaderLen,
		TotalLen: ipv4.HeaderLen + len(b),
		// ID:       12345,
		Protocol: int(proto),
		// Src:      net.ParseIP("192.168.25.1").To4(),
		Dst: net.ParseIP("127.0.0.1").To4(),
		TTL: 64,
	}

	if err := iprawconn.WriteTo(h, b, nil); err != nil {
		log.Println(err)
	}
	log.Println("iprawconn.WriteTo done")

	buf := make([]byte, 1024)
	for {
		_, p, _, err := iprawconn.ReadFrom(buf)
		if err != nil {
			fmt.Printf("Error reading: %#v\n", err)
			return
		}
		fmt.Printf("Message received: %s\n", string(p))
		fmt.Printf("Message received: %s\n", string(buf))
	}
}
