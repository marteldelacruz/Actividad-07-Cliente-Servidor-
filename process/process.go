package process

import (
	"fmt"
	"time"
)

type Task interface {
	RunProcess()
}

type ProcessAdmin struct {
	Processes []Task
}

type Process struct {
	PrintValues      bool
	TerminateProcess bool
	I                uint64
	ID               uint64
}

func (process *Process) RunProcess() {
	for {
		// print values
		if process.PrintValues {
			fmt.Printf("id %d: %d\n", process.ID, process.I)
		}

		process.I = process.I + 1
		time.Sleep(time.Millisecond * 500)
		// terminate process
		if process.TerminateProcess {
			break
		}
	}
}
