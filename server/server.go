package server

import (
	"fmt"
	"sync"
	"syscall"
	"tcp-server/tcp"
	"tcp-server/client_t"
)

var connMutex sync.Mutex
var isn = 5000
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

	fmt.Printf("Server is listening on port %d\n",s.port)

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
	var isn = 5000
	client := client.Client_t{
		Fd: fd,
		Port: h.Src_port,
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

	synAckPacket := tcp.ContructTCPHeader(h.Dest_port, h.Src_port, uint32(isn), h.Seq+1, uint8(5), uint8(0), tcp.TCP_SYNACK, uint16(0), uint16(0), uint16(0))
	if err := syscall.Sendto(fd, synAckPacket, 0, &destAddr); err != nil {
		panic("failed to send SYN-ACK packet " + err.Error())
	}
	isn++
	fmt.Println("sent syn ack")
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

		retHeader := tcp.ParseTCPHeader(buf[20:])
		if retHeader.Src_port == h.Src_port && retHeader.Ack == uint32(0x1001)+1 && retHeader.Flags == tcp.TCP_ACK {
			connMutex.Lock()
				s.activeConnections[h.Src_port] = true
			connMutex.Unlock()
			fmt.Println("Connection established with", from)
			/*connMutex.Lock()
			delete(activeConnections, h.Src_port)
			connMutex.Unlock()*/
		}
		if s.activeConnections[retHeader.Src_port]&& retHeader.Src_port == client.Port && retHeader.Flags==tcp.TCP_PSH{
			msg := string(buf[40:])	
			fmt.Printf("%d: %s",retHeader.Src_port,msg)
			ackPacket := tcp.ContructTCPHeader(s.port, h.Src_port, uint32(isn), uint32(h.Seq)+uint32(n-40), uint8(5), uint8(0), tcp.TCP_SYNACK, uint16(0), uint16(0), uint16(0))
			if err := syscall.Sendto(fd, ackPacket, 0, &destAddr); err != nil {
				panic("failed to send SYN-ACK packet " + err.Error())
			}
		}
	}
}


