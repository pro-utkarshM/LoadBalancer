// worker.go

package main

import (
	"log"
	"net"
	"net/rpc"
	"os/exec"
)

// BuildTask represents a compilation task.
type BuildTask struct {
	Language string // "C", "C++", "Objective-C", etc.
	Source   string // Source code file or directory path
}

// BuildResponse represents the result of a build task.
type BuildResponse struct {
	Output string // Compilation output or errors
}

// Worker performs compilation tasks.
type Worker struct{}

// PerformBuild compiles the given source code.
func (w *Worker) PerformBuild(task BuildTask, response *BuildResponse) error {
	cmd := exec.Command("gcc", task.Source) // Change "gcc" to the appropriate compiler
	output, err := cmd.CombinedOutput()
	if err != nil {
		response.Output = string(output) + "\nBuild failed."
	} else {
		response.Output = string(output) + "\nBuild successful."
	}
	return nil
}

func main() {
	// Create a worker
	worker := new(Worker)

	// Register worker as an RPC server
	rpc.Register(worker)

	// Listen for incoming RPC requests
	l, err := net.Listen("tcp", ":9001") // Specify the worker's address
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	log.Println("Worker listening on port 9001...")
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go rpc.ServeConn(conn)
	}
}
