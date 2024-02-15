package dbopts

import (
	"database/sql"
	"fmt"
)

func Insert(db *sql.DB, user, script, next_run string) error {

	err := Create(db, user)
	if err != nil {
		return err
	}

	id := setID(db, user)
	status := "Done"

	cmd := fmt.Sprintf("INSERT INTO %s (id, script, next_run, status) VALUES (%d, '%s', '%s', '%s')", user, id, script, next_run, status)
	fmt.Println(cmd)
	_, err = db.Exec(cmd)
	if err != nil {
		return err
	}

	return err

}

func setID(db *sql.DB, table string) int {
	cmd := "SELECT MAX(id) AS max_id FROM " + table
	var maxID int
	err := db.QueryRow(cmd).Scan(&maxID)
	if err != nil {
		if maxID == 0 {
			return 1
		}
	}
	maxID += 1

	return maxID
}
