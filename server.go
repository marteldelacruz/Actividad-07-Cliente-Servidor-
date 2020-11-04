package main

import (
	"encoding/gob"
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
	// creates the new processes
	startProcesses(processAdmin)

	// loop to handle clients
	for {
		client, err := server.Accept()

		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleClient(client, processAdmin)
	}
}

// This function takes charge of handling clients
// by sending them a process the first time they
// connect to the server
func handleClient(client net.Conn, processAdmin *process.ProcessAdmin) {
	lastProcess := processAdmin.Processes[len(processAdmin.Processes)-1]
	// sending process
	err := gob.NewEncoder(client).Encode(lastProcess)

	// terminate when an error ocurrs
	if err != nil {
		fmt.Println(err)
		return
	} else {
		// remove last element
		processAdmin.Processes = processAdmin.Processes[:len(processAdmin.Processes)-1]
	}
}

func main() {
	processAdmin := process.ProcessAdmin{}
	go server(&processAdmin)
	// press enter to exit...
	fmt.Scanln()
}
