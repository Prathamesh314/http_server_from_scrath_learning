## reading from channel and printing the whole line by reading 8 bytes at a time
```
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
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

	filepath := "message.txt"
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Error in opening file: %s\n", err.Error())
	}

	lines := getLinesChannel(f)

	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}
}
```
