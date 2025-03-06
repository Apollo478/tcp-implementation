package main

import (
	"fmt"
	"encoding/binary"
	"syscall"
	"tcp-server/tcp"
)
func main(){
	fd,err :=  syscall.Socket(syscall.AF_INET,syscall.SOCK_RAW,syscall.IPPROTO_TCP)
	if err != nil {
	 	panic("could not create file descriptor " +err.Error())	
	}
	defer syscall.Close(fd)
	if err := syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_HDRINCL, 1); err != nil {
		panic("Failed to set IP_HDRINCL: " + err.Error())
	}
	dstAddr := syscall.SockaddrInet4{
		Port: 8080,
		Addr: [4]byte{127,0,0,1},
	}
	synPacket := buildPacket();
	if err := syscall.Sendto(fd,synPacket,0,&dstAddr);err!=nil{
		panic("unable to send packet "+err.Error())
	}
	fmt.Println("packet sent")
}
func buildPacket() []byte{
	ipHeader := []byte{
		0x45, 0x00, 0x28, 0x00, 
		0x00, 0x01, 0x00, 0x00, 
		0x40, 0x06, 0x00, 0x00, 
		0x7F, 0x00, 0x00, 0x01, 	
		0x7F, 0x00, 0x00, 0x01,
	}

	tcpHeader := tcp.TCPheader{
		Src_port: 8081,
		Dest_port: 8080,
		Seq: 0x1000,
		Ack: 0,
		Hlen: 5,
		Flags: syscall.TCP_SYNCNT,
		Window: 0x1000,
	}
	packet := append(ipHeader, serializeTCPHeader(&tcpHeader)...);
	fmt.Println(len(packet))
	return packet;
}
func serializeTCPHeader(h *tcp.TCPheader) []byte{
	data := make([]byte,20)
	binary.BigEndian.PutUint16(data[0:2], h.Src_port)
    binary.BigEndian.PutUint16(data[2:4], h.Dest_port)
    binary.BigEndian.PutUint32(data[4:8], h.Seq)
    binary.BigEndian.PutUint32(data[8:12], h.Ack)
    data[12] = (h.Hlen << 4) | (h.Res << 1)
    data[13] = h.Flags
    binary.BigEndian.PutUint16(data[14:16], h.Window)
    binary.BigEndian.PutUint16(data[16:18], h.Checksum)
    binary.BigEndian.PutUint16(data[18:20], h.Urg)
    return data
}



