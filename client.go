package main

import (
	"fmt"
	"net"
)

func client() {
	conn, err := net.Dial("tcp", ":9999")

	if err != nil {
		fmt.Println(err)
		return
	}

	msg := "Hola lola"
	conn.Write([]byte(msg))
	conn.Close()
}

func main() {
	go client()

	fmt.Scanln()
}
