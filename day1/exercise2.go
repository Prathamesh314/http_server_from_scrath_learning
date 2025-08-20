package day1

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func readFromFile1(filepath string) error {
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}

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
			fmt.Printf("read: %s\n", ans)
			ans = ""
		}

		ans += string(data)
	}

	if len(ans) != 0 {
		fmt.Printf("read: %s\n", ans)
	}

	return nil
}