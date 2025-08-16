package main

import (
	"io"
	"os"
)

func readFromFile1(filepath string) (string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return "", err
	}

	var ans string

	for {
		buffer := make([]byte, 8)
		n1, err := f.Read(buffer)

		if err == io.EOF{
			return ans, nil
		}
		if err != nil {
			return "", err
		}

		ans += string(buffer[:n1])
	}
}