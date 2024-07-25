package service

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"ont/internal/dbopts"
	"ont/internal/escape"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/exp/rand"
)

type Message struct {
	Command string        `json:"command"`
	User    string        `json:"user"`
	Job     dbopts.Jobs   `json:"job"`
	Jobs    []dbopts.Jobs `json:"jobs"`
	Status  string        `json:"status"`
}

func Server(db *sql.DB, ip, port string) {

	serverAddr := fmt.Sprintf("%s:%s", ip, port)
	listener, err := net.Listen("tcp", serverAddr)
	if err != nil {
		escape.LogFatal(err.Error())
	}

	defer listener.Close()
	escape.LogPrint("Ontd server running on port 3033")

	for {

		conn, _ := listener.Accept()
		reader := bufio.NewReader(conn)
		message, err := reader.ReadString('\n')
		if err != nil {
			escape.LogPrint("Error reading from connection:", err)
		}

		var msg Message

		err = json.Unmarshal([]byte(message), &msg)
		if err != nil {
			escape.LogPrint(err.Error())
		}

		//Verbose logging
		//escape.LogPrintf("User '%s' requested '%s' job \n", msg.User, msg.Command)

		var response Message
		fun := []string{"Okay.", "Cool.", "Roger.", "Got it.", "On it.", "Sure.", "All right.", "Certainly.", "Will do.", "Absolutely."}

		switch msg.Command {
		case "list":
			jobs, err := dbopts.List(db, msg.User)
			if err != nil {
				escape.LogPrint(err.Error())
			}

			response := Message{
				Command: msg.Command,
				User:    msg.User,
				Jobs:    jobs,
			}

			//Verbose logging
			//escape.LogPrint(jobs)

			sendResponse(response, conn)

		case "run":
			jobID, err := msg.Job.Insert(db, msg.User, true)
			if err != nil {
				escape.LogPrint(err.Error())
			}

			n := rand.Intn(len(fun))

			status := fmt.Sprintf("%s New job '%d' is created.", fun[n], jobID)
			response := Message{
				Status: status,
			}
			sendResponse(response, conn)

		case "stop":
			err := msg.Job.ChangeJobStatus(db, msg.User, "Inactive", false)
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
				n := rand.Intn(len(fun))
				status := fmt.Sprintf("%s Job %d is inactive now.", fun[n], msg.Job.Id)
				response = Message{
					Status: status,
				}
			}
			sendResponse(response, conn)

		case "start":
			err := msg.Job.ChangeJobStatus(db, msg.User, "Active", false)
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
				n := rand.Intn(len(fun))
				status := fmt.Sprintf("%s Job %d is active now.", fun[n], msg.Job.Id)
				response = Message{
					Status: status,
				}
			}
			sendResponse(response, conn)

		case "refresh":
			job, err := msg.Job.GetJob(db, msg.User, msg.Job.Id)
			if err != nil {
				escape.LogPrint(err.Error())
			}

			if job.Status == "Inactive" {
				response = Message{
					Status: "Can't refresh an inactive job.",
				}
			} else {
				err = msg.Job.ChangeJobStatus(db, msg.User, "Active", true)
				if err != nil {
					fmt.Println(err.Error())
					if err.Error() == "sql: no rows in result set" {
						status := fmt.Sprintf("Job %d doesn't exist", msg.Job.Id)
						response = Message{
							Status: status,
						}
					}
				} else {
					n := rand.Intn(len(fun))
					status := fmt.Sprintf("%s Job %d is refreshed.", fun[n], msg.Job.Id)
					response = Message{
						Status: status,
					}
				}
			}
			sendResponse(response, conn)

		case "remove":
			//if err := dbopts.RemoveJob(db, msg.User, msg.Job); err != nil {
			if err := msg.Job.RemoveJob(db, msg.User); err != nil {
				escape.LogPrint(err.Error())
				status := fmt.Sprintf("Job %d doesn't exist.", msg.Job.Id)
				response = Message{
					Status: status,
				}
			} else {
				n := rand.Intn(len(fun))
				status := fmt.Sprintf("%s Job %d is removed.", fun[n], msg.Job.Id)
				response = Message{
					Status: status,
				}
			}
			sendResponse(response, conn)

		case "clean":
			status := "All your entries are removed"
			err := dbopts.CleanAllJobs(db, msg.User)
			if err != nil {
				//Verbose Logging
				if strings.Contains(err.Error(), "Unknown table") {
					status = "You don't have any jobs to clean."
				} else {
					escape.LogPrint(err.Error())
				}
			}

			response := Message{
				Status: status,
			}
			sendResponse(response, conn)
		}
		conn.Close()
	}

}

func sendResponse(response any, conn net.Conn) {

	responseData, err := json.Marshal(response)
	if err != nil {
		escape.LogPrint("Error marshaling JSON:", err)
	}

	_, err = conn.Write(append(responseData, '\n'))
	if err != nil {
		escape.LogPrint("Error writing to UDP:", err)
	}
}
