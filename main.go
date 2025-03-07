package main

import (
	"fmt"
	"tcp-server/server"
)
func main(){
	fmt.Println("hello there");
	s,err := server.NewServer(8080)
	if err != nil{
		fmt.Println(err.Error())
	}
	s.Listen()
}
