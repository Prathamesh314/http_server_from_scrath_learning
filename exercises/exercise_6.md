## Replaced getLinesFromChannel with our RequestReader
```
package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/Prathamesh314/http_server_from_scrath_learning/internal/request"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	out_ch := make(chan string, 8)

	go func() {
		defer close(out_ch)
		ans := ""
		for {
			buffer := make([]byte, 8)

			n, err := f.Read(buffer)
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Error in reading into buffer: %s\n", err.Error())
				break
			}

			data := buffer[:n]
			if idx := strings.Index(string(data), "\n"); idx != -1 {
				ans += string(data)[:idx]
				data = data[idx+1:]
				out_ch <- ans
				ans = ""
			}

			ans += string(data)
		}

		if len(ans) != 0 {
			out_ch <- ans
		}
	} ()
	return out_ch
}

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

		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", r.RequestLine.Method, r.RequestLine.RequestTarget, r.RequestLine.HttpVersion)	
	}
}
```
