##### To read 8 bytes at a time but print whole line not just 8 bytes

```
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)


func main() {
	// now you have to read data in same way
	// read data 8 bytes at a time.
	// but now print the entire line. not just 8 bytes

	filepath := "message.txt"
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Error in opening file: %s\n", err.Error())
	}

	// create an empty string for newline
	ans := ""
	for {
		// create a buffer of 8 bytes
		buffer := make([]byte, 8)
		n, err := f.Read(buffer)

		if err == io.EOF{
			break
		}

		if err != nil {
			log.Fatalf("Error in reading data into buffer: %s\n", err.Error())
			return
		}

		data := buffer[:n]
		// now we have to find \n in data and if it is present then append that much part to ans print it then update the ans
		if idx := strings.Index(string(data), "\n"); idx != -1 {
			ans += string(data)[:idx]
			data = data[idx+1: ]
			fmt.Printf("read: %s\n", ans)
			ans = ""
		}

		ans += string(data)
	}

	if len(ans) != 0 {
		fmt.Printf("read: %s\n", ans)
	}
}
```