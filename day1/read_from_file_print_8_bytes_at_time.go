package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func readFromFile(filepath string) error{
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}

	for {
		b := make([]byte, 8)
		n1, err := f.Read(b)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading message.txt: %s", err)
		}

		fmt.Printf("read: %s\n", b[:n1])
	}
}

func main(){
	// fmt.Printf("I hope I get the job!")
	err := readFromFile("message.txt")
	if err != nil {
		log.Fatalf("Error while reading message.txt: %s\n", err)
		return
	}
}