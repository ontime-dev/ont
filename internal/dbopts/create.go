package dbopts

import (
	"database/sql"
	"fmt"

	"github.com/ontime-dev/ont/internal/escape"

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

func Create(db *sql.DB, user string, verbose bool) error {

	if verbose {
		escape.LogPrintf("DEBUG(CRT): Creating new table for user %s", user)
	}

	//cmd := "CREATE TABLE IF NOT EXISTS " + user + table + ";"
	cmd := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s %s", user, table)
	_, err := db.Exec(cmd)

	if err != nil {
		escape.LogPrint(err.Error())
		return err
	}

	//verbose logging
	if verbose {
		escape.LogPrintf("DEBUG(CRT): Table %s created", user)
	}

	return err
}
