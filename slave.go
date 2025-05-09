package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Read the command from master
	command, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading command:", err)
		return
	}

	// Strip newline
	command = command[:len(command)-1]

	if command == "shutdown" {
		// Perform shutdown on the slave machine
		fmt.Println("💻 Shutting down...")
		err := exec.Command("shutdown", "/s", "/f", "/t", "0").Run()
		if err != nil {
			fmt.Println("Error shutting down:", err)
		}
	} else if command == "restart" {
		// Handle restart command
		fmt.Println("💻 Restarting...")
		err := exec.Command("shutdown", "/r", "/f", "/t", "0").Run()
		if err != nil {
			fmt.Println("Error restarting:", err)
		}
	}
}

func main() {
	// Listen for incoming connections
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("❌ Failed to start server:", err)
		os.Exit(1)
	}
	defer listener.Close()

	// Notify master that this slave is online
	go notifyMaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func notifyMaster() {
	// Notify the master that this slave is online
	conn, err := net.Dial("tcp", "192.168.1.3:8080") // Master IP
	if err != nil {
		fmt.Println("❌ Failed to connect to master:", err)
		return
	}
	defer conn.Close()

	fmt.Fprintf(conn, "online\n")
}
