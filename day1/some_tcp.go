package day1

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)


func readFromFile4(ch chan string, f io.ReadCloser) {
	ans := ""

	for {
		buffer := make([]byte, 8)
		n1, err := f.Read(buffer)
		if err == io.EOF{
			break
		}
		data := buffer[:n1]
		if i := bytes.IndexByte(data, '\n'); i != -1{
			ans += string(data[:i])
			data = data[i+1:]
			// fmt.Printf("read: %s\n", ans)
			ch <- ans
			ans = ""
		}

		ans += string(data)
	}

	if len(ans) != 0 {
		ch <- ans
	}

	ch <- "EOF"

}

func getLinesChannel2(f io.ReadCloser) <-chan string {
	ch := make(chan string, 0)

	go readFromFile4(ch, f)

	return ch
}


func main(){
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("Error while listening on port 42069: %s", err.Error())
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Error while accepting the connection: %s", err.Error())
		}

		lines := getLinesChannel2(conn)

		for line := range lines {
			if line == "EOF" {
				break
			}

			fmt.Printf("read: %s\n", line)
		}
		// fmt.Printf("Connection received: %s", conn.LocalAddr())
	}
}