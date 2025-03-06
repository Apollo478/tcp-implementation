package tcp

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
