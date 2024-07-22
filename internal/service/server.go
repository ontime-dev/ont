package service

import (
	"bufio"
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

func Server(db *sql.DB, port string) {

	/*addr := net.UDPAddr{
		Port: 3033,
		IP:   net.ParseIP("127.0.0.1"),
	}*/

	//conn, err := net.ListenUDP("udp", &addr)
	portNum := fmt.Sprintf(":%s", port)
	listener, err := net.Listen("tcp", portNum)
	if err != nil {
		escape.LogFatal(err.Error())
	}

	defer listener.Close()
	escape.LogPrint("Ontd server running on port 3033")

	//buffer := make([]byte, 1024)

	for {
		//n, clientAddr, err := conn.ReadFromUDP(buffer)
		conn, _ := listener.Accept()
		reader := bufio.NewReader(conn)
		message, err := reader.ReadString('\n')
		if err != nil {
			escape.LogPrint("Error reading from connection:", err)
		}

		var msg Message
		//defer conn.Close()
		//err = json.Unmarshal(buffer[:n], &msg)
		err = json.Unmarshal([]byte(message), &msg)
		if err != nil {
			escape.LogPrint(err.Error())
		}

		escape.LogPrintf("User '%s' requested '%s' job \n", msg.User, msg.Command)

		var response Message

		switch msg.Command {
		case "list":
			err, jobs := dbopts.List(db, msg.User)
			if err != nil {
				escape.LogPrint(err.Error())
			}

			response := Message{
				Command: msg.Command,
				User:    msg.User,
				Jobs:    jobs,
			}

			sendResponse(response, conn)

		case "run":
			err := dbopts.Insert(db, msg.User, msg.Job, true)
			if err != nil {
				escape.LogPrint(err.Error())
			}
			response := Message{
				Status: "Ok",
			}
			sendResponse(response, conn)

		case "stop":
			err := dbopts.ChangeJobStatus(db, msg.User, "Inactive", msg.Job)
			if err != nil {
				if err.Error() == "sql: no rows in result set" {
					status := fmt.Sprintf("Job %d doesn't exist", msg.Job.Id)
					response = Message{
						Status: status,
					}
				} else {
					escape.LogPrintf(err.Error())
					status := fmt.Sprintf("Job %d is already inactive.", msg.Job.Id)
					response = Message{
						Status: status,
					}
				}
			} else {
				status := fmt.Sprintf("Job %d is inactive now.", msg.Job.Id)
				response = Message{
					Status: status,
				}
			}
			sendResponse(response, conn)

		case "start":
			err := dbopts.ChangeJobStatus(db, msg.User, "Active", msg.Job)
			if err != nil {
				if err.Error() == "sql: no rows in result set" {
					status := fmt.Sprintf("Job %d doesn't exist", msg.Job.Id)
					response = Message{
						Status: status,
					}
				} else {
					escape.LogPrintf(err.Error())
					status := fmt.Sprintf("Job %d is already active.", msg.Job.Id)
					response = Message{
						Status: status,
					}
				}
			} else {
				status := fmt.Sprintf("Job %d is active now.", msg.Job.Id)
				response = Message{
					Status: status,
				}
			}
			sendResponse(response, conn)

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
			sendResponse(response, conn)
		}

		//	:= dbopts.Opt(msg.Command, msg.User, msg.Job, cfgFile)
		conn.Close()
	}

}

// func sendResponse(response any, clientAddr *net.UDPAddr, conn *net.UDPConn) {
func sendResponse(response any, conn net.Conn) {

	responseData, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
	}

	//_, err = conn.WriteToUDP(responseData, clientAddr)
	_, err = conn.Write(append(responseData, '\n'))
	if err != nil {
		fmt.Println("Error writing to UDP:", err)
	}
}
