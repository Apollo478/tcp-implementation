package tcp

import (
	"fmt"
	"syscall"
)

func Listen(){
	fd,err := syscall.Socket(syscall.AF_INET,syscall.SOCK_RAW,syscall.IPPROTO_TCP)
	if err!=nil {
		panic("failed to create file descriptor "+err.Error())
	}
	defer syscall.Close(fd)
	if err := syscall.SetsockoptInt(fd,syscall.IPPROTO_IP,syscall.IP_HDRINCL,1) ;err != nil{
		panic("failed to set ip header include "+err.Error())
	}
	fmt.Println("server is listening on port 8080")
	for{
		buf := make([]byte,1500)
		n,from,err := syscall.Recvfrom(fd,buf,0)
		if err!=nil{
			fmt.Printf("Error recieving message %s\n",err.Error());
		}
		if n < 40 {
			fmt.Println("ip and tcp header shorter than 40")
		}
		tcpHeader := ParseTCPHeader(buf[20:]);
		if tcpHeader.Dest_port==8080 {
			
			fmt.Printf("recieved tcp packet from %v\n",from);
			fmt.Printf("source port %d\n",tcpHeader.Src_port);
			fmt.Printf("destination port %d\n",tcpHeader.Dest_port);
			fmt.Printf("flags  %d\n",tcpHeader.Flags);
			fmt.Printf("data size is %d",len(buf))
		}
	}
}

