package dbopts

import (
	"database/sql"

	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Jobs struct {
	Id       int
	Script   string
	Next_run string
	Every    string
	Status   string
}

func Opt(operation, user string, job Jobs) error {

	db, err := sql.Open("mysql", "ont:password@/ontime")
	if err != nil {
		return err
	}
	defer db.Close()
	fmt.Printf("Operation: %s -- User: %s\n", operation, user)

	switch operation {
	case "create":
		fmt.Println("Operation: create")
		err = Create(db, user)
		if err != nil {
			return err
		}
	case "insert":
		//arg[1] = user, arg[2] = script path, arg[3] = next_run
		fmt.Println("Hello there")
		err := Insert(db, user, job)
		if err != nil {
			return err
		}
	case "list":
		//arg[1] = user
		err := FetchJobs(db, user)
		if err != nil {
			return err
		}
	}

	return nil
}
