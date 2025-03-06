package main

import (
	"fmt"
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
	synPacket := tcp.ContructTCPHeader(uint16(8081),uint16(8080),uint32(0x1000),uint32(0),uint8(5),uint8(0),syscall.TCP_SYNCNT,uint16(0),uint16(0),uint16(0))
	if err := syscall.Sendto(fd,synPacket,0,&dstAddr);err!=nil{
		panic("unable to send packet "+err.Error())
	}
	fmt.Println("packet sent")
}



