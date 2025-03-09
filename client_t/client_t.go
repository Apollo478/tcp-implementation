package client

import (
	"bufio"
	"fmt"
	"os"
	"syscall"
	"tcp-server/tcp"
	"time"
)
var serverPort = 42069
type Client_t struct {
	Fd int
	Port uint16
}
func NewClient(port uint16) *Client_t{
	fd,err :=  syscall.Socket(syscall.AF_INET,syscall.SOCK_RAW,syscall.IPPROTO_TCP)
	if err != nil {
	 	panic("could not create file descriptor " +err.Error())	
	}
	if err := syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_HDRINCL, 1); err != nil {
		panic("Failed to set IP_HDRINCL: " + err.Error())
	}
	client := &Client_t{
		Fd : fd,
		Port : port,
	}
	return client
}
func (c *Client_t) Start(){
	c.connect()
}
func (c *Client_t) connect(){
	dstAddr := syscall.SockaddrInet4{
		Port: serverPort,
		Addr: [4]byte{127,0,0,1},
	}
	var isn = 1000
	start := time.Now()
	synPacket := tcp.ContructTCPHeader(uint16(c.Port),uint16(serverPort),uint32(isn),uint32(0),uint8(5),uint8(0),tcp.TCP_SYN,uint16(0),uint16(0),uint16(0))
	if err := syscall.Sendto(c.Fd,synPacket,0,&dstAddr);err!=nil{
		panic("unable to send packet "+err.Error())
	}
	fmt.Println("packet sent")
	isn++
	for {
		if time.Since(start) > 5*time.Second {
			panic("Timeout waiting for server response")
		}
		buf := make([]byte,1500)
		n,from,err := syscall.Recvfrom(c.Fd,buf,0)
		if err!=nil{
			fmt.Printf("Error recieving message %s\n",err.Error());
		}
		if n < 40 {
			fmt.Println("ip and tcp header shorter than 40")
		}
		tcpHeader := tcp.ParseTCPHeader(buf[20:]);
		if tcpHeader.Src_port == uint16(serverPort){
			start = time.Now()
			fmt.Printf("recieved packet from %v",from)
			if tcpHeader.Flags == tcp.TCP_PSH {
				msgBytes := buf[40:]
				msg := string(msgBytes)
				fmt.Printf("message is %s \n",msg)
			}
			if tcpHeader.Flags == tcp.TCP_ACK {
				if tcpHeader.Ack != uint32(isn){
					fmt.Println("invalid ack")
				}
			}
			if tcpHeader.Flags == tcp.TCP_SYNACK  && tcpHeader.Ack == uint32(isn){
				fmt.Println("recieved connection reply from server")	
				synPacket := tcp.ContructTCPHeader(tcpHeader.Dest_port,tcpHeader.Src_port,uint32(isn),tcpHeader.Seq+1,uint8(5),uint8(0),tcp.TCP_ACK,uint16(0),uint16(0),uint16(0))
				err := syscall.Sendto(c.Fd,synPacket,0,&dstAddr)
				if err!= nil {
					panic("failed to send packet "+err.Error())
				}
				par := uint32(isn)
				go c.communicate(&par)
			}
		}
	}
}
func (c *Client_t) communicate(seq *uint32){
	reader := bufio.NewReader(os.Stdin)
	for{
		fmt.Print("Enter message : ")
		msg,err := reader.ReadString('\n')
		if err!=nil{
			panic("fialed to read message from input"+err.Error())
		}
		if err := tcp.Send(uint16(c.Port), uint16(serverPort), c.Fd, msg); err != nil {
			panic("Failed to send message: " + err.Error())
		}
		*seq +=uint32(len([]byte(msg)))
	}
}
