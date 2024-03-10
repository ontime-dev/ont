package dbopts

import (
	"database/sql"
	"errors"

	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Jobs struct {
	Id        int
	Script    string
	Exec_time string
	Every     string
	Status    string
}

func Opt(operation, user string, job Jobs) error {

	db, err := sql.Open("mysql", "ont:password@/ontime")
	if err != nil {
		return err
	}
	defer db.Close()
	//fmt.Printf("Operation: %s -- User: %s\n", operation, user)

	switch operation {
	case "create":
		fmt.Println("Operation: create")
		err = Create(db, user)
		if err != nil {
			return err
		}
	case "insert":
		//arg[1] = user, arg[2] = script path, arg[3] = next_run
		//fmt.Println("Hello there")
		err := Insert(db, user, job, true)
		if err != nil {
			return err
		}
	case "list":
		//arg[1] = user
		if job.Id == 0 {
			err := PrintJobs(db, user)
			if err != nil {
				return err
			}
		} else if job.Id > 0 {
			err := PrintOneJob(db, user, job.Id)
			if err != nil {
				return err
			}
		} else {
			return errors.New("please enter a valid jobID")
		}
	case "stop":
		err := ChangeJobStatus(db, user, "Inactive", job)
		if err != nil {
			return err
		}
	case "start":
		err := ChangeJobStatus(db, user, "Active", job)
		if err != nil {
			return err
		}
	case "remove":
		if err := RemoveJob(db, user, job); err != nil {
			return err
		}

	}

	return nil
}
