package day1

import (
	"bytes"
	"io"
)

func readFromFile3(ch chan string, f io.ReadCloser) {
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

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string, 0)

	go readFromFile3(ch, f)

	return ch
}
