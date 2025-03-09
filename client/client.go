package main

import (
	"os"
	"strconv"
	"tcp-server/client_t"
)
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
	client := client.NewClient(uint16(port))
	client.Start()
}
