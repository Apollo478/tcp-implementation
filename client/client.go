package client

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"syscall"
	"tcp-server/tcp"
)
type Client_t struct {
	Fd int
	Port uint16
}
var serverPort = 42069
func main(){
	args := os.Args
	if len(args)<2{
		panic("No argument provided");
	}
	port,err := strconv.ParseInt(args[1],10,16) 
	if err!=nil {
		panic("failed to parse arguments"+err.Error())
	}
	Client(int(port))
}
func Client(port int){
	
	fd,err :=  syscall.Socket(syscall.AF_INET,syscall.SOCK_RAW,syscall.IPPROTO_TCP)
	if err != nil {
	 	panic("could not create file descriptor " +err.Error())	
	}
	defer syscall.Close(fd)
	if err := syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_HDRINCL, 1); err != nil {
		panic("Failed to set IP_HDRINCL: " + err.Error())
	}
	dstAddr := syscall.SockaddrInet4{
		Port: serverPort,
		Addr: [4]byte{127,0,0,1},
	}
	synPacket := tcp.ContructTCPHeader(uint16(port),uint16(serverPort),uint32(0x1000),uint32(0),uint8(5),uint8(0),tcp.TCP_SYN,uint16(0),uint16(0),uint16(0))
	if err := syscall.Sendto(fd,synPacket,0,&dstAddr);err!=nil{
		panic("unable to send packet "+err.Error())
	}
	fmt.Println("packet sent")
	for {
		
		buf := make([]byte,1500)
		n,from,err := syscall.Recvfrom(fd,buf,0)
		if err!=nil{
			fmt.Printf("Error recieving message %s\n",err.Error());
		}
		if n < 40 {
			fmt.Println("ip and tcp header shorter than 40")
		}
		tcpHeader := tcp.ParseTCPHeader(buf[20:]);
		if tcpHeader.Src_port == uint16(serverPort){
			fmt.Printf("recieved packet from %v",from)
			if tcpHeader.Flags == tcp.TCP_PSH {
				msgBytes := buf[40:]
				msg := string(msgBytes)
				fmt.Printf("message is %s \n",msg)
			}
			if tcpHeader.Flags == tcp.TCP_SYNACK  && tcpHeader.Ack == uint32(0x1000)+1{
				fmt.Println("recieved connection reply from server")	
				synPacket := tcp.ContructTCPHeader(tcpHeader.Dest_port,tcpHeader.Src_port,uint32(0x1002),tcpHeader.Seq+1,uint8(5),uint8(0),tcp.TCP_ACK,uint16(0),uint16(0),uint16(0))
				err := syscall.Sendto(fd,synPacket,0,&dstAddr)
				if err!= nil {
					panic("failed to send packet "+err.Error())
				}
				break;
			}
		}
	}
	reader := bufio.NewReader(os.Stdin)
	for{
		fmt.Print("Enter message : ")
		msg,err := reader.ReadString('\n')
		if err!=nil{
			panic("fialed to read message from input"+err.Error())
		}
		if err := tcp.Send(uint16(port), uint16(serverPort), fd, msg); err != nil {
			panic("Failed to send message: " + err.Error())
		}
	}
}


