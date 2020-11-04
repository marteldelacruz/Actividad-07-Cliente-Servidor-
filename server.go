package main

import (
	"fmt"
	"net"

	process "./process"
)

var PROTOCOL = "tcp"
var PORT = ":9999"
var TOTAL_PROCESSES = 5

// Adds the initial processes count to the admin
// This also set the params to start printing values
func startProcesses(processAdmin *process.ProcessAdmin) {
	for i := 0; i < TOTAL_PROCESSES; i++ {
		// brand new process
		proc := process.Process{
			PrintValues:      true,
			TerminateProcess: false,
			I:                0,
			ID:               uint64(i + 1),
		}
		processAdmin.Processes = append(processAdmin.Processes, &proc)
		// start go routine
		go processAdmin.Processes[i].RunProcess()
	}
}

/// This rutine runs the server on a loop to keep
/// handling client petitions using the TCP connection
/// on the 9999 port
func server(processAdmin *process.ProcessAdmin) {
	server, err := net.Listen(PROTOCOL, PORT)

	// terminate when an error ocurrs
	if err != nil {
		fmt.Println(err)
		return
	}

	startProcesses(processAdmin)
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

	// terminate when an error ocurrs
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("Message: ", data[:br])
	}
}

func main() {
	processAdmin := process.ProcessAdmin{}
	go server(&processAdmin)

	//process.Process()
	fmt.Scanln()
}
