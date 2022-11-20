package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		//go handleConn(conn)
		Msg(conn)
	}
}
func handleConn(c net.Conn) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05\n\r"))
		if err != nil {
			return
		}
		time.Sleep(1 * time.Second)

	}
}
func Msg(b net.Conn) {
	defer b.Close()

	for {
		var a string
		fmt.Print("Введите сообщение: ")
		_, err := fmt.Scanln(&a)
		if err != nil {
			fmt.Println("Некорректный ввод", err)
			continue
		}

		if n, err := b.Write([]byte(a)); n == 0 || err != nil {
			fmt.Println(err)
			return
		}
	}
}
