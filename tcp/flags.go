package tcp
const (
	TCP_FIN  = 0x01
	TCP_SYN  = 0x02
	TCP_RST  = 0x04
	TCP_PSH  = 0x08
	TCP_ACK  = 0x10
	TCP_URG  = 0x20

	TCP_SYNACK = TCP_SYN | TCP_ACK
)
