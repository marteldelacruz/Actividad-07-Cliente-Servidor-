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
			PrintValues: true,
			Terminate:   false,
			I:           0,
			ID:          uint64(i + 1),
		}
		processAdmin.Processes = append(processAdmin.Processes, &proc)
		// start go routine
		go processAdmin.Processes[i].RunProcess()
	}
}

// Verifies if a client adress already exist on the list
func clientIsInList(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

/// This rutine runs the server on a loop to keep
/// handling client petitions using the TCP connection
/// on the 9999 port
func server(processAdmin *process.ProcessAdmin) {
	server, err := net.Listen(PROTOCOL, PORT)
	var adressList []string

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

		go handleClient(client, processAdmin, &adressList)
	}
}

// This function takes charge of handling clients
// by sending them a process the first time they
// connect to the server
func handleClient(client net.Conn, processAdmin *process.ProcessAdmin, clientsList *[]string) {
	// wait for client ID
	data := make([]byte, 100)
	br, err := client.Read(data)

	// terminate when an error ocurrs
	if err != nil {
		fmt.Println(err)
		return
	} else {
		clientID := string(data[:br])
		// check if client already exist
		if clientIsInList(clientID, *clientsList) {
			recieveProcessFromClient(client, processAdmin)
		} else {
			*clientsList = append(*clientsList, clientID)
			sendProcessToClient(client, processAdmin)
		}
	}

}

func sendProcessToClient(client net.Conn, processAdmin *process.ProcessAdmin) {
	lastIndex := len(processAdmin.Processes) - 1
	lastProcess := processAdmin.Processes[lastIndex]

	// stop last process
	processAdmin.Processes[lastIndex].StopProcess()

	// sending process
	err := gob.NewEncoder(client).Encode(lastProcess)

	// terminate when an error ocurrs
	if err != nil {
		fmt.Println(err)
		return
	} else {
		// remove last process from list
		processAdmin.Processes = processAdmin.Processes[:lastIndex]
	}
}

func recieveProcessFromClient(client net.Conn, processAdmin *process.ProcessAdmin) {
	process := process.Process{}

	// wait for process to be returned back
	err := gob.NewDecoder(client).Decode(&process)

	// terminate when an error ocurrs
	if err != nil {
		fmt.Println(err)
		return
	} else {
		// run process
		process.ContinueProcess()
		process.RunProcess()
		processAdmin.Processes = append(processAdmin.Processes, &process)
	}
}

func main() {
	processAdmin := process.ProcessAdmin{}
	go server(&processAdmin)

	// press enter to exit...
	fmt.Scanln()
}
