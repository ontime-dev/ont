package dbopts

import (
	"database/sql"

	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func Opt(args ...string) error {

	db, err := sql.Open("mysql", "ont:password@/ontime")
	if err != nil {
		return err
	}
	defer db.Close()
	fmt.Printf("Operation: %s -- User: %s\n", args[0], args[1])

	switch args[0] {
	case "create":
		fmt.Println("Operation: create")
		err = Create(db, args[1])
		if err != nil {
			return err
		}
	case "insert":
		//arg[1] = user, arg[2] = script path
		//Insert(db, args[1], args[2])
	}

	return nil
}
