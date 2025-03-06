package tcp

import (
	"fmt"
	"syscall"
	"encoding/binary"
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
		tcpHeader := parseTCPHeader(buf[20:]);
		if tcpHeader.Dest_port==8080 {
			
			fmt.Printf("recieved tcp packet from %v\n",from);
			fmt.Println("haffy youre a big dummy")
			/*fmt.Printf("source port %d\n",tcpHeader.Src_port);
			fmt.Printf("destination port %d\n",tcpHeader.Dest_port);
			fmt.Printf("flags  %08b\n",tcpHeader.Flags);
			fmt.Println(syscall.TCP_SYNCNT)*/
		}
	}

}

func parseTCPHeader(data []byte) *TCPheader {
	if len(data) < 20{
		return nil
	}
	header := &TCPheader{
		Src_port: binary.BigEndian.Uint16(data[0:2]),
		Dest_port: binary.BigEndian.Uint16(data[2:4]),
		Seq: binary.BigEndian.Uint32(data[4:8]),
		Ack: binary.BigEndian.Uint32(data[8:12]),
		Hlen: uint8((data[12] & 0xF0)>> 4), 
		Res: uint8((data[12] & 0x0F)>> 1), 
		Flags: uint8((data[13] & 0x3F)), 
		Window: binary.BigEndian.Uint16(data[14:16]),
		Checksum: binary.BigEndian.Uint16(data[16:18]),
		Urg: binary.BigEndian.Uint16(data[18:20]),
	}
	return header
}
