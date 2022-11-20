package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	buf := make([]byte, 512)
	for {
		_, err = conn.Read(buf)
		if err == io.EOF {
			break
		}
		io.WriteString(os.Stdout, fmt.Sprintf("Custom output! %s", string(buf)))
	}
}
