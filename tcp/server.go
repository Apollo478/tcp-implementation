package tcp

import (
	"fmt"
	"sync"
	"syscall"
)

var activeConnections = make(map[uint16]bool)
var connMutex sync.Mutex

func Listen() {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_TCP)
	if err != nil {
		panic("failed to create file descriptor " + err.Error())
	}
	defer syscall.Close(fd)

	if err := syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_HDRINCL, 1); err != nil {
		panic("failed to set IP_HDRINCL " + err.Error())
	}

	fmt.Println("Server is listening on port 8080")

	for {
		buf := make([]byte, 1500)
		n, from, err := syscall.Recvfrom(fd, buf, 0)
		if err != nil {
			fmt.Printf("Error receiving message: %s\n", err.Error())
			continue
		}
		if n < 40 {
			fmt.Println("IP and TCP header shorter than 40 bytes")
			continue
		}

		tcpHeader := ParseTCPHeader(buf[20:])
		if tcpHeader.Dest_port == 8080 {
			connMutex.Lock()
			if activeConnections[tcpHeader.Src_port] {
				connMutex.Unlock()
				continue
			}
			connMutex.Unlock()

			fmt.Printf("Received TCP packet from %v\n", from)
			fmt.Printf("Source Port: %d\n", tcpHeader.Src_port)
			fmt.Printf("Destination Port: %d\n", tcpHeader.Dest_port)
			fmt.Printf("Flags: %d\n", tcpHeader.Flags)
			fmt.Printf("Data Size: %d\n", len(buf))

			if tcpHeader.Flags == TCP_SYN {
				connMutex.Lock()
				activeConnections[tcpHeader.Src_port] = true
				connMutex.Unlock()

				go handleConnection(*tcpHeader)
			}
		}
	}
}

func handleConnection(h TCPheader) {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_TCP)
	if err != nil {
		panic("failed to create file descriptor " + err.Error())
	}
	defer syscall.Close(fd)

	if err := syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_HDRINCL, 1); err != nil {
		panic("failed to set IP_HDRINCL " + err.Error())
	}

	fmt.Println("Received connection request")

	destAddr := syscall.SockaddrInet4{
		Port: int(h.Src_port),
		Addr: [4]byte{127, 0, 0, 1},
	}

	synAckPacket := ContructTCPHeader(h.Dest_port, h.Src_port, uint32(0x1001), h.Seq+1, uint8(5), uint8(0), TCP_SYNACK, uint16(0), uint16(0), uint16(0))
	if err := syscall.Sendto(fd, synAckPacket, 0, &destAddr); err != nil {
		panic("failed to send SYN-ACK packet " + err.Error())
	}

	for {
		buf := make([]byte, 1500)
		n, from, err := syscall.Recvfrom(fd, buf, 0)
		if err != nil {
			fmt.Printf("Error receiving message: %s\n", err.Error())
			continue
		}
		if n < 40 {
			fmt.Println("IP and TCP header shorter than 40 bytes")
			continue
		}

		ackHeader := ParseTCPHeader(buf[20:])
		if ackHeader.Src_port == h.Src_port && ackHeader.Ack == uint32(0x1001)+1 && ackHeader.Flags == TCP_ACK {
			fmt.Println("Connection established with", from)
			/*connMutex.Lock()
			delete(activeConnections, h.Src_port)
			connMutex.Unlock()*/
			break
		}
	}
}

