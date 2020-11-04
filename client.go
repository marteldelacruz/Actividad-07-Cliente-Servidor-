package main

import (
	"encoding/gob"
	"fmt"
	"net"

	process "./process"
)

var _PROTOCOL = "tcp"
var _PORT = ":9999"

// The new client will connect to a server on the PORT
// if the server is not running already, an error message
// will be shown
func client(process *process.Process) {
	conn, err := net.Dial(_PROTOCOL, _PORT)

	if err != nil {
		fmt.Println(err)
		return
	}

	handleServerProcess(conn, process)
}

// Handles the process recieved from the server
func handleServerProcess(conn net.Conn, process *process.Process) {
	err := gob.NewDecoder(conn).Decode(process)

	// terminate when an error ocurrs
	if err != nil {
		fmt.Println(err)
		return
	} else {
		// run process
		process.TerminateProcess = false
		process.RunProcess()
	}
}

func main() {
	process := process.Process{}
	go client(&process)

	fmt.Scanln()
}
