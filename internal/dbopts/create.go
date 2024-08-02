package dbopts

import (
	"database/sql"
	"fmt"

	"github.com/ontime-dev/ont/internal/escape"

	_ "github.com/go-sql-driver/mysql"
)

const table = `(
	timestamp TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3),
	id INT,
	script VARCHAR(255),
	exec_time VARCHAR(255),
	every VARCHAR(255),
	status VARCHAR(255),
	runon VARCHAR(255)
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
