package client

import (
	"encoding/json"
	"fmt"
	"net"
	"ont/internal/dbopts"
)

type Message struct {
	Command string        `json:"command"`
	User    string        `json:"user"`
	Job     []dbopts.Jobs `json:"job"`
}

func SendMsg(message any) error {
	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:3033")
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return err
	}

	// Dial the server address
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Println("Error dialing:", err)
		return err
	}
	defer conn.Close()

	messageData, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	_, err = conn.Write(messageData)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println("Error receiving response:", err)
	}
	fmt.Printf("Raw buffer content: %s\n", string(buffer[:n]))

	var response Message
	err = json.Unmarshal(buffer[:n], &response)
	fmt.Println(buffer[:n])
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
	}

	fmt.Println("Server response:", response.Job)

	return nil
}

/*
func RecieveRspns() error {
	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:3033")
	if err != nil {
		escape.LogPrint("Error resolving address:", err)
		return err
	}

	// Dial the server address
	conn, err := net.DialUDP("udp", nil, serverAddr)
	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		escape.LogPrint("Error receiving response:", err)
	}

	defer conn.Close()

	var response Message
	err = json.Unmarshal(buffer[:n], &response)
	if err != nil {
		escape.LogPrint("Error unmarshaling JSON:", err)
	}

	escape.LogPrint("Server response:", response)

	return nil
}*/
