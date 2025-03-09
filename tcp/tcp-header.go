package tcp

import (
	"encoding/binary"
	"fmt"
	"syscall"
)

type TCPheader struct {
	Src_port uint16
	Dest_port uint16
	Seq uint32
	Ack uint32
	Hlen uint8
	Res uint8
	Flags uint8
	Window uint16
	Checksum uint16
	Urg uint16
}
func ContructTCPHeader(Src_port uint16,Dest_port uint16,Seq uint32,Ack uint32,Hlen uint8,Res uint8,Flags uint8,Window uint16,Checksum uint16,Urg uint16) []byte {
	
	ipHeader := []byte{
		0x45, 0x00, 0x28, 0x00, 
		0x00, 0x01, 0x00, 0x00, 
		0x40, 0x06, 0x00, 0x00, 
		0x7F, 0x00, 0x00, 0x01, 	
		0x7F, 0x00, 0x00, 0x01,
	}

	tcpHeader := TCPheader{
		Src_port:	Src_port,
		Dest_port:	Dest_port,
		Seq:	Seq,
		Ack:	Ack,
		Hlen:	Hlen,
		Flags:	Flags,
		Window:	Window,
	}
	packet := append(ipHeader, SerializeTCPHeader(&tcpHeader)...);
	return packet;
}

func SerializeTCPHeader(h *TCPheader) []byte{
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

func ParseTCPHeader(data []byte) *TCPheader {
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
func  Send(srcPort uint16,destPort  uint16,fd int,msg string) error{
	seq:= /*rand.Int() % math.MaxInt32*/ 12
	msgPacket := ContructTCPHeader(srcPort, destPort, uint32(seq), 0, uint8(5), uint8(0), TCP_PSH, uint16(0), uint16(0), uint16(0))
	msgBytes := []byte(msg)
	packet := append(msgPacket,msgBytes...)
	fmt.Println(destPort,srcPort)
	destAddr := syscall.SockaddrInet4{
		Port: int(destPort),
		Addr: [4]byte{127, 0, 0, 1},
	}

	if err := syscall.Sendto(fd, packet, 0, &destAddr); err != nil {
		return err
	}
	return nil;
}
