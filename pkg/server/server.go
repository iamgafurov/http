package server

import (
	"net"
	"sync"
	"io"
	"log"
	"strings"
	"bytes"
	"errors"
)

type HandlerFunc func(conn net.Conn)

type Server struct {
	addr string
	mu sync.RWMutex
	handlers map[string]HandlerFunc
}

func NewServer(addr string)*Server{
	return &Server{addr:addr,handlers: make(map[string]HandlerFunc)}
}

func (s *Server) Register(path string, handler HandlerFunc){
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[path]= handler
}

func (s *Server) Start() error{
	listener, err := net.Listen("tcp",s.addr)
	if err != nil {
		log.Print(err)
		return err
	}
	
	


	for {
		conn,err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go s.handle(conn)
		if err != nil {
			log.Print(err)
			continue
		}
	}

	return nil
}

func (s *Server)handle (conn net.Conn)(err error){
	defer func(){
		if cerr := conn.Close();cerr !=nil {
			if err == nil{
				err = cerr
				return
			}
		}
		log.Print(err)
	}()
	
	buf := make ([]byte,4096)
	n,err := conn.Read(buf)
	if err == io.EOF {
		return errors.New("request is empty")
	}
	if err != nil {
		return err
	}

	data := buf[:n]
	requestLineDlim := []byte{'\r','\n'}
	requestLineEnd := bytes.Index(data,requestLineDlim)
	if requestLineEnd == -1 {
		return nil
	}

	requestLine := string(data[:requestLineEnd])
	parts := strings.Split(requestLine," ")
	if len(parts) != 3 {
		return errors.New("parts len is not 3")
	}

	path := parts[1]
	s.mu.RLock()
	for name,handler := range(s.handlers){
		if name == path {
			s.mu.RUnlock()
			handler(conn)
			break
		}
	}
	return nil	
}
