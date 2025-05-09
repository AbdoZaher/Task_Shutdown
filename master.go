package main

import (
	"fmt"
	"net"
	"time"


)

var slaveIPs = []string{
	"192.168.1.8", // Add the IPs of the slave devices here
	"192.168.1.3",
}

func sendCommand(ip, command string) {
	conn, err := net.Dial("tcp", ip+":8080")
	if err != nil {
		fmt.Println("❌ Failed to connect to", ip)
		return
	}
	defer conn.Close()

	fmt.Fprintf(conn, "%s\n", command)
	fmt.Println("✅ Command sent successfully to", ip)
}

func ping(ip string) bool {
	// Just a simple check for connectivity
	conn, err := net.Dial("tcp", ip+":8080")
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func main() {
	for {
		// Scan for active slaves
		activeSlaves := []string{}
		for _, ip := range slaveIPs {
			if ping(ip) {
				activeSlaves = append(activeSlaves, ip)
			}
		}

		// Show active slaves to the user
		fmt.Println("Active Slaves:")
		for _, ip := range activeSlaves {
			fmt.Println(ip)
		}

		// Ask master to choose a slave to shutdown
		var slaveIP string
		fmt.Println("Enter IP of slave to shutdown:")
		fmt.Scanln(&slaveIP)

		if ping(slaveIP) {
			sendCommand(slaveIP, "shutdown")
		} else {
			fmt.Println("❌ Slave not found or already offline!")
		}

		// Wait for a while before checking again
		time.Sleep(5 * time.Second)
	}
}
