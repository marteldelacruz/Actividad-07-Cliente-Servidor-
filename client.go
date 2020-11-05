package main

import (
	"encoding/gob"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"

	process "./process"
)

var _PROTOCOL = "tcp"
var _PORT = ":9999"

// The new client will connect to a server on the PORT
// if the server is not running already, an error message
// will be shown
func client(process *process.Process, id string) {
	conn, err := net.Dial(_PROTOCOL, _PORT)

	if err != nil {
		fmt.Println(err)
		return
	}
	// send client ID
	conn.Write([]byte(id))

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
		conn.Close()
	}
}

// Returns the process to the server
func returnProcessToServer(process *process.Process, id string) {
	conn, err := net.Dial(_PROTOCOL, _PORT)

	if err != nil {
		fmt.Println(err)
		return
	}
	// send client ID
	conn.Write([]byte(id))

	err = gob.NewEncoder(conn).Encode(*process)

	// terminate when an error ocurrs
	if err != nil {
		fmt.Println(err)
		return
	} else {
		process.StopProcess()
	}

}

func generateID() string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	return strconv.Itoa(r1.Intn(100))
}

func main() {
	process := process.Process{}
	id := generateID()

	go client(&process, id)
	fmt.Scanln()

	// return process to server
	go returnProcessToServer(&process, id)

	fmt.Scanln()
}
