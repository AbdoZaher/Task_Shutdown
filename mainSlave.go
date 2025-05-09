package main

import (
	"bufio"
	//"fmt"
	"log"
	"net"
	"os/exec"
	"runtime"
)

func main() {
	listener, err := net.Listen("tcp", ":9090") // Listening on port 9090
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer listener.Close()
	log.Println("Slave is running. Waiting for master...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Connection error:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	cmd, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Failed to read command:", err)
		return
	}

	log.Printf("Received command: %s", cmd)

	switch cmd {
	case "shutdown\n":
		log.Println("Shutdown command received")
		shutdown()
	default:
		log.Println("Unknown command")
	}
}

func shutdown() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("shutdown", "/s", "/f", "/t", "0")
	} else {
		cmd = exec.Command("shutdown", "-h", "now")
	}

	err := cmd.Run()
	if err != nil {
		log.Println("Failed to shutdown:", err)
	} else {
		log.Println("Shutdown executed")
	}
}
