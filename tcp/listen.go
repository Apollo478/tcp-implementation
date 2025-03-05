package tcp

import (
	"fmt"
	"syscall"
	"encoding/binary"
)

func Listen(port int16){
	fd,err := syscall.Socket(syscall.AF_INET,syscall.SOCK_RAW,syscall.IPPROTO_TCP)
	if err!=nil {
		panic("failed to create file descriptor "+err.Error())
	}
	defer syscall.Close(fd)
	if err := syscall.SetsockoptInt(fd,syscall.IPPROTO_IP,syscall.IP_HDRINCL,1) ;err != nil{
		panic("failed to set ip header include "+err.Error())
	}
	fmt.Println("server is listening on port ")
	for{
		buf := make([]byte,1500)
		n,from,err := syscall.Recvfrom(fd,buf,0)
		if err!=nil{
			fmt.Printf("Error recieving message %s\n",err.Error());
		}
		if n < 40 {
			fmt.Println("ip and tcp header shorter than 40")
		}
		tcpHeader := parseTCPHeader(buf[20:]);
		fmt.Printf("recieved tcp packet from %v\n",from);
		fmt.Printf("source port %d\n",tcpHeader.src_port);
		fmt.Printf("destination port %d\n",tcpHeader.src_dest);
		fmt.Printf("flags  %08b\n",tcpHeader.flags);
		fmt.Println(syscall.TCP_SYNCNT)
	}

}

func parseTCPHeader(data []byte) *TCPheader {
	if len(data) < 20{
		return nil
	}
	header := &TCPheader{
		src_port: binary.BigEndian.Uint16(data[0:2]),
		src_dest: binary.BigEndian.Uint16(data[2:4]),
		seq: binary.BigEndian.Uint32(data[4:8]),
		ack: binary.BigEndian.Uint32(data[8:12]),
		hlen: uint8((data[12] & 0xF0)>> 4), 
		res: uint8((data[12] & 0x0F)>> 1), 
		flags: uint8((data[13] & 0x3F)), 
		window: binary.BigEndian.Uint16(data[14:16]),
		checksum: binary.BigEndian.Uint16(data[16:18]),
		urg: binary.BigEndian.Uint16(data[18:20]),
	}
	return header
}
