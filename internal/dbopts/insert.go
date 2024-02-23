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
	every := "Daily"

	cmd := fmt.Sprintf("INSERT INTO %s (id, script, next_run, every, status) VALUES (%d, '%s', '%s', '%s', '%s')", user, id, script, next_run, every, status)
	fmt.Println(cmd)
	_, err = db.Exec(cmd)
	if err != nil {
		fmt.Printf("HEREEEEEE")
		return err
	}

	return err

}

func setID(db *sql.DB, table string) int {
	cmd := fmt.Sprintf("SELECT MAX(id) AS max_id FROM %s", table)
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
