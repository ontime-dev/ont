package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"ont/internal/dbopts"
)

type Message struct {
	Command string        `json:"command"`
	User    string        `json:"user"`
	Job     dbopts.Jobs   `json:"job"`
	Jobs    []dbopts.Jobs `json:"jobs"`
	Status  string        `json:"status"`
}

func SendMsg(message Message) (error, Message) {
	// serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:3033")
	// if err != nil {
	// 	fmt.Println("Error resolving address:", err)
	// 	return err, Message{}
	// }

	// Dial the server address
	conn, err := net.Dial("tcp", "127.0.0.1:3033")
	if err != nil {
		fmt.Println("Error dialing:", err)
		return err, Message{}
	}
	defer conn.Close()

	messageData, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err, Message{}
	}

	_, err = conn.Write(append(messageData, '\n'))
	if err != nil {
		fmt.Println(err.Error())
		return err, Message{}
	}

	// buffer := make([]byte, 2048)
	// n, _, err := conn.ReadFromUDP(buffer)
	// if err != nil {
	// 	fmt.Println("Error receiving response:", err)
	// }
	//fmt.Printf("Raw buffer content: %s\n", string(buffer[:n]))
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading from connection:", err)
	}

	var responsemsg Message
	err = json.Unmarshal([]byte(response), &responsemsg)
	// fmt.Println(n)
	// fmt.Println(buffer)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
	}

	//fmt.Println("Server response:", response.Job)

	return nil, responsemsg
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
