package server

import (
	"fmt"
	"sync"
	"syscall"
	"tcp-server/tcp"
)

var connMutex sync.Mutex
type Server struct {
	fd int
	port uint16
	activeConnections map[uint16]bool
}
func NewServer(port uint16) (*Server, error){
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_TCP)
	if err != nil {
		return nil, fmt.Errorf("failed to create socket %v",err)
	}

	if err := syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_HDRINCL, 1); err != nil {
		syscall.Close(fd)
		return nil, fmt.Errorf("failed to set ip header inclusion %v",err)
	}
	return &Server{fd: fd, port:port, activeConnections:  make(map[uint16]bool)} ,nil
}
func (s *Server)Listen() {

	fmt.Printf("Server is listening on port %d",s.port)

	for {
		buf := make([]byte, 1500)
		n, from, err := syscall.Recvfrom(s.fd, buf, 0)
		if err != nil {
			fmt.Printf("Error receiving message: %s\n", err.Error())
			continue
		}
		if n < 40 {
			fmt.Println("IP and TCP header shorter than 40 bytes")
			continue
		}

		tcpHeader := tcp.ParseTCPHeader(buf[20:])
		if tcpHeader.Dest_port == s.port {
			connMutex.Lock()
			if s.activeConnections[tcpHeader.Src_port] {
				connMutex.Unlock()
				continue
			}
			connMutex.Unlock()

			fmt.Printf("Received TCP packet from %v\n", from)
			fmt.Printf("Source Port: %d\n", tcpHeader.Src_port)
			fmt.Printf("Destination Port: %d\n", tcpHeader.Dest_port)
			fmt.Printf("Flags: %d\n", tcpHeader.Flags)
			fmt.Printf("Data Size: %d\n", len(buf))

			if tcpHeader.Flags == tcp.TCP_SYN {

				go s.handleConnection(*tcpHeader)
			}
		}
	}
}

func  (s *Server) handleConnection(h tcp.TCPheader) {
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

	synAckPacket := tcp.ContructTCPHeader(h.Dest_port, h.Src_port, uint32(0x1001), h.Seq+1, uint8(5), uint8(0), tcp.TCP_SYNACK, uint16(0), uint16(0), uint16(0))
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

		ackHeader := tcp.ParseTCPHeader(buf[20:])
		if ackHeader.Src_port == h.Src_port && ackHeader.Ack == uint32(0x1001)+1 && ackHeader.Flags == tcp.TCP_ACK {
			s.activeConnections[h.Src_port] = true
			connMutex.Lock()
			connMutex.Unlock()
			fmt.Println("Connection established with", from)
			/*connMutex.Lock()
			delete(activeConnections, h.Src_port)
			connMutex.Unlock()*/
			break
		}
	}
}

