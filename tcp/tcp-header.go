package tcp

type TCPheader struct {
	src_port uint16;
	src_dest uint16;
	seq uint32;
	ack uint32;
	hlen uint8;
	res uint8;
	flags uint8;
	window uint16;
	checksum uint16;
	urg uint16;
}
