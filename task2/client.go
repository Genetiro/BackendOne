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
	var name string
	fmt.Print("Введите имя: ")
	fmt.Scanln(&name)

	// отправляем сообщение серверу
	if n, err := conn.Write([]byte(name)); n == 0 || err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		io.Copy(os.Stdout, conn)
	}()
	io.Copy(conn, os.Stdin)
	fmt.Printf("%s: exit", conn.LocalAddr())
}
