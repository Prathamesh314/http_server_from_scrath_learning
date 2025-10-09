package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Prathamesh314/http_server_from_scrath_learning/internal/request"
)

func main() {
	// now you have to read data in same way
	// read data 8 bytes at a time.
	// but now print the entire line. not just 8 bytes

	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("Error while listening tcp at port 42069: %s\n", err.Error())
		return
	}
	// log.Printf("Listening at port 42069.....")

	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Error in accepting new connection: %s\n", err.Error())
			return
		}

		r, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("error while Reading request: %s\n", err.Error())
			return
		}

		// fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", r.RequestLine.Method, r.RequestLine.RequestTarget, r.RequestLine.HttpVersion)	
		fmt.Printf("Request line:\n")
		fmt.Printf("- Method: %s\n- Target: %s\n- Version: %s\n", r.RequestLine.Method, r.RequestLine.RequestTarget, r.RequestLine.HttpVersion)
		fmt.Println("Headers: ")
		for key, val := range r.Headers {
			fmt.Printf("- %s: %s\n", key, val)
		}
		fmt.Printf("Body:\n %s\n", r.Body)
	}
}
