##### To read 8 bytes at a time.
```
package main

import (
	"fmt"
	"io"
	"log"
	"os"
)


func main() {
	// fmt.Printf("I hope I get the job!");
	// Read data from file
	// you have to read 8 bytes per time from the file

	// so let's create a buffer of 8 bytes
	
	// open a file
	filename := "message.txt"
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error in opening file: %s\n", err.Error())
		return
	}
	
	for {
		// create buffer of 8 bytes
		buffer := make([]byte, 8);
		n, err := f.Read(buffer)

		if err == io.EOF {{
			return
		}}

		if err != nil {
			log.Fatalf("Error reading data to buffer: %s", err.Error())
		}
		// get 8 bytes data
		data := buffer[:n]
		// print it
		fmt.Printf("read: %s\n", data)
	}
}
```