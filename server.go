package main

import (
	"fmt"
	"net"
)

/// This rutine runs the server on a loop to keep
/// handling client petitions using the TCP connection
/// on the 9999 port
func server() {
	server, err := net.Listen("tcp", ":9999")

	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		client, err := server.Accept()

		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleClient(client)
	}
}

func handleClient(client net.Conn) {
	data := make([]byte, 100)
	br, err := client.Read(data)

	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("Message: ", data[:br])
	}
}

func main() {
	go server()

	fmt.Scanln()
}
