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

// type Message struct {
// 	Command string        `json:"command"`
// 	User    string        `json:"user"`
// 	Job     []dbopts.Jobs `json:"job"`
// }

type Message struct {
	Command string        `json:"command"`
	User    string        `json:"user"`
	Job     dbopts.Jobs   `json:"job"`
	Jobs    []dbopts.Jobs `json:"jobs"`
	Status  string        `json:"status"`
}

func Server(db *sql.DB) {

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

		escape.LogPrintf("User '%s' requested '%s' job \n", msg.User, msg.Command)

		var response Message

		switch msg.Command {
		case "list":
			err, jobs := dbopts.List(db, msg.User)
			if err != nil {
				escape.Error(err.Error())
			}

			response := Message{
				Command: msg.Command,
				User:    msg.User,
				Jobs:    jobs,
			}

			sendResponse(response, clientAddr, conn)

		case "run":
			err := dbopts.Insert(db, msg.User, msg.Job, true)
			if err != nil {
				escape.Error(err.Error())
			}
			response := Message{
				Status: "Ok",
			}
			sendResponse(response, clientAddr, conn)

		case "stop":
			err := dbopts.ChangeJobStatus(db, msg.User, "Inactive", msg.Job)
			if err != nil {
				escape.LogPrintf(err.Error())
				status := fmt.Sprintf("Job %d is already inactive.", msg.Job.Id)
				response = Message{
					Status: status,
				}
			} else {
				status := fmt.Sprintf("Job %d is inactive now.", msg.Job.Id)
				response = Message{
					Status: status,
				}
			}
			sendResponse(response, clientAddr, conn)

		case "start":
			err := dbopts.ChangeJobStatus(db, msg.User, "Active", msg.Job)
			if err != nil {
				escape.LogPrintf(err.Error())
				status := fmt.Sprintf("Job %d is already active.", msg.Job.Id)
				response = Message{
					Status: status,
				}
			} else {
				status := fmt.Sprintf("Job %d is active now.", msg.Job.Id)
				response = Message{
					Status: status,
				}
			}
			sendResponse(response, clientAddr, conn)

		case "remove":

			if err := dbopts.RemoveJob(db, msg.User, msg.Job); err != nil {
				escape.LogPrint(err.Error())
				status := fmt.Sprintf("Job %d doesn't exist.", msg.Job.Id)
				response = Message{
					Status: status,
				}
			} else {
				status := fmt.Sprintf("Job %d is removed.", msg.Job.Id)
				response = Message{
					Status: status,
				}
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
