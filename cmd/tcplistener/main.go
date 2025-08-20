package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer f.Close()
		defer close(out)
		var ans bytes.Buffer
		for {
			data := make([]byte, 8)
			n, err := f.Read(data)
			if err != nil {
				if err == io.EOF {
					// Send whatever remains in the buffer at EOF
					if ans.Len() > 0 {
						log.Printf("Ans: %s\n", ans.String())
						out <- ans.String()
					}
				}
				break
			}

			data = data[:n]

			if i := bytes.IndexByte(data, '\n'); i != -1 {
				ans.Write(data[:i])
				out <- ans.String()
				ans.Reset()
				ans.Write(data[i+1:])
			} else {
				ans.Write(data)
			}
		}
	}()
	return out
}

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("Error while listening on port 42069")
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Error while accepting connection on port 42069")
		}

		for line := range getLinesChannel(conn) {
			fmt.Printf("read: %s\n", line)
		}
	}
}