package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"ont/internal/config"
	"ont/internal/dbopts"
	"ont/internal/escape"
)

type Message struct {
	Command string        `json:"command"`
	User    string        `json:"user"`
	Job     dbopts.Jobs   `json:"job"`
	Jobs    []dbopts.Jobs `json:"jobs"`
	Status  string        `json:"status"`
}

func SendMsg(message Message) (Message, error) {

	// Dial the server address
	port := config.GetConfig("SERVER_PORT")
	serverAddr := fmt.Sprintf("%s:%s", "localhost", port)

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		escape.Error(err.Error())
		return Message{}, err
	}
	defer conn.Close()

	messageData, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return Message{}, err
	}

	_, err = conn.Write(append(messageData, '\n'))
	if err != nil {
		fmt.Println(err.Error())
		return Message{}, err
	}

	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		escape.Error("Error reading from connection:", err)
	}

	var responsemsg Message
	err = json.Unmarshal([]byte(response), &responsemsg)

	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
	}

	return responsemsg, nil
}
