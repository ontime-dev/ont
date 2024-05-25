package service

import (
	"encoding/json"
	"net"
	"ont/internal/escape"
)

type Message struct {
	Command string `json:"command"`
	User    string `json:"user"`
}

func Server() {
	addr := net.UDPAddr{
		Port: 3033,
		IP:   net.ParseIP("127.0.0.1"),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		escape.Error(err.Error())
	}

	defer conn.Close()
	escape.LogPrint("Ontd server running on port 3033")

	buffer := make([]byte, 1024)

	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			escape.Error(err.Error())
		}

		var msg Message
		err = json.Unmarshal(buffer[:n], &msg)
		if err != nil {
			escape.Error(err.Error())
		}

		//DO WHAT SHOULD BE DONE WITH THE JSON
		escape.LogPrintf("Received message from %s: %s\n", clientAddr, msg.Command)
	}

}
