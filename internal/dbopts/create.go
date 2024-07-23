package dbopts

import (
	"database/sql"
	"ont/internal/escape"

	_ "github.com/go-sql-driver/mysql"
)

const table = `(
	timestamp timestamp(3),
	id int,
	script varchar(255),
	exec_time varchar(255),
	every varchar(255),
	status varchar(255),
	runon varchar(255)
)`

func Create(db *sql.DB, user string) error {

	cmd := "CREATE TABLE IF NOT EXISTS " + user + table + ";"
	//verbose logging
	//escape.LogPrintf("Table %s created", user)
	_, err := db.Exec(cmd)

	if err != nil {
		escape.LogPrint(err.Error())
		return err
	}

	return nil

}
