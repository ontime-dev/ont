package dbopts

import (
	"database/sql"
	"fmt"
)

func CleanAllJobs(db *sql.DB, table string) error {
	cmd := fmt.Sprintf("Drop table %s;", table)

	_, err := db.Exec(cmd)
	if err != nil {
		return err
	}

	return nil

}
