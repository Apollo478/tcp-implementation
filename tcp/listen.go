package tcp

import (
	"fmt"
	"syscall"
)

func Listen(port int16){
	fd,err := syscall.Socket(syscall.AF_INET,syscall.SOCK_RAW,syscall.IPPROTO_TCP)
	if err!=nil {
		panic("failed to create file descriptor "+err.Error())
	}
	defer syscall.Close(fd)
}
