package server

import (
	"log"
	"net"
	"strconv"

	"github.com/Prathamesh314/http_server_from_scrath_learning/internal/response"
)

type Server struct {
	listener net.Listener
	port int
}

func (s *Server) Close() error {
	return s.listener.Close()
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	
	body := "Hello World!"
	headers := response.GetDefaultHeaders(len(body))
	
	if err := response.WriteStatusLine(conn, response.StatusCodeOK); err != nil {
		log.Printf("Failed to write status line: %v", err)
		return
	}
	
	if err := response.WriteHeaders(conn, headers); err != nil {
		log.Printf("Failed to write headers: %v", err)
		return
	}
	
	if _, err := conn.Write([]byte(body)); err != nil {
		log.Printf("Failed to write body: %v", err)
	}
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return
		}

		log.Printf("Accepting new connection from address: %s\n", conn.RemoteAddr())
		go s.handle(conn)
	}
}

func Serve(port int) (*Server, error) {
	log.Printf("Started http server on port %d!", port)
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}

	server := Server{
		listener: listener,
		port: port,
	}

	go server.listen()

	return &server, nil
}