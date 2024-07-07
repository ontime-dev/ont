package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"ont/internal/dbopts"
	"ont/internal/escape"

	_ "github.com/go-sql-driver/mysql"
)

type Message struct {
	Command string        `json:"command"`
	User    string        `json:"user"`
	Job     []dbopts.Jobs `json:"job"`
}

func Server() {
	//cfgFile := "/etc/ont/ont.conf"
	db, err := sql.Open("mysql", "ont:password@/ontime")
	if err != nil {
		escape.Error(err.Error())
	}
	defer db.Close()

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
		escape.LogPrintf("Received message from %s: %s, %s \n", clientAddr, msg.Command, msg.User)

		switch msg.Command {
		case "list":
			/*type Response struct {
				Jobs []dbopts.Jobs `json:"jobs"`
			}*/
			err, jobs := dbopts.List(db, msg.Job, msg.User)
			if err != nil {
				escape.Error(err.Error())
			}
			//response := Response{Jobs: jobs}
			response := Message{
				Command: msg.Command,
				User:    msg.User,
				Job:     jobs,
			}
			sendResponse(response, clientAddr, conn)

		}
		//	:= dbopts.Opt(msg.Command, msg.User, msg.Job, cfgFile)

	}

}

func sendResponse(response any, clientAddr *net.UDPAddr, conn *net.UDPConn) {
	responseData, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
	}

	_, err = conn.WriteToUDP(responseData, clientAddr)
	if err != nil {
		fmt.Println("Error writing to UDP:", err)
	}
}
