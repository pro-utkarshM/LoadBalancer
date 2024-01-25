// coordinator.go

package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
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

// BuildCoordinator manages build tasks and coordinates with workers.
type BuildCoordinator struct {
	Workers []string // Worker addresses
}

// Build performs a distributed build.
func (coordinator *BuildCoordinator) Build(task BuildTask, response *BuildResponse) error {
	// Simple round-robin task distribution to workers
	workerAddr := coordinator.Workers[len(coordinator.Workers)%len(coordinator.Workers)]
	client, err := rpc.Dial("tcp", workerAddr)
	if err != nil {
		return err
	}
	defer client.Close()

	// Forward the build task to the worker
	err = client.Call("Worker.PerformBuild", task, response)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// Set up the coordinator
	coordinator := &BuildCoordinator{
		Workers: []string{"localhost:9001", "localhost:9002"}, // Add worker addresses
	}

	// Register coordinator as an RPC server
	rpc.Register(coordinator)

	// Listen for incoming RPC requests
	l, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	fmt.Println("Coordinator listening on port 9000...")
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go rpc.ServeConn(conn)
	}
}

