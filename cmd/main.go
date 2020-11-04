package main

import (
	"http/pkg/server"
	"os"
	"net"
	"strconv"
	"log"
)
func main(){
	host := "0.0.0.0"
	port :="9999"

	if err := execute(host,port); err !=nil {
		os.Exit(1)
	}
}

func execute (host string , port string)(err error){
	srv := server.NewServer(net.JoinHostPort(host, port))
	srv.Register("/",func(conn net.Conn){
		body:= "welcome to our web-site"
		_,err = conn.Write(generateResponse(body))
		if err != nil {
			log.Print(err)
		}
	})
	srv.Register("/about",func(conn net.Conn){
		body:= "About Golang Academy"
		_,err = conn.Write(generateResponse(body))
		if err != nil {
			log.Print(err)
		}
	})
	log.Print("server run in ",host +":"+ port)
	return srv.Start()
}

func generateResponse(body string)[]byte{
	return ([]byte("HTTP/1.1 200 OK\r\n"+
	"Content-Length: "+ strconv.Itoa(len(body)) + "\r\n"+
	"Content-Type: text/html\r\n"+
	"Connection: close\r\n"+
	"\r\n" + 
	body,))
}