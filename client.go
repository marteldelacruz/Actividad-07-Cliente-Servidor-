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
func client(process *process.Process, _conn *net.Conn) {
	conn, err := net.Dial(_PROTOCOL, _PORT)
	*_conn = conn

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
		process.Terminate = false
		process.RunProcess()
	}
}

// Returns the process to the server
func returnProcessToServer(process *process.Process, conn *net.Conn) {
	err := gob.NewEncoder(*conn).Encode(*process)

	// terminate when an error ocurrs
	if err != nil {
		fmt.Println(err)
		return
	} else {
		process.StopProcess()
	}

}

func main() {
	process := process.Process{}
	var conn net.Conn

	go client(&process, &conn)
	fmt.Scanln()

	// return process to server
	go returnProcessToServer(&process, &conn)

	fmt.Scanln()
	conn.Close()
}
