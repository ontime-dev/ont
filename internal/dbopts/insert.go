package dbopts

import (
	"database/sql"
	"fmt"
)

func Insert(db *sql.DB, user string, job Jobs) error {

	err := Create(db, user)
	if err != nil {
		return err
	}

	id := setID(db, user)
	status := "Done"

	cmd := fmt.Sprintf("INSERT INTO %s (id, script, exec_time, every, status) VALUES (%d, '%s', '%s', '%s', '%s');", user, id, job.Script, job.Exec_time, job.Every, status)
	//cmd := fmt.Sprintf("INSERT INTO %s (script, next_run, every, status) VALUES ( '%s', '%s', '%s', '%s')", user, job.Script, job.Next_run, job.Every, status)

	//fmt.Println(cmd)
	_, err = db.Exec(cmd)
	if err != nil {
		return err
	}

	return err

}

func setID(db *sql.DB, table string) int {
	/*
		cmd := fmt.Sprintf("SELECT MAX(id) AS max_id FROM %s", table)
		var maxID int
		err := db.QueryRow(cmd).Scan(&maxID)
	*/
	maxID, err := getMaxID(db, table)
	if err != nil {
		if maxID == 0 {
			return 1
		}
	}
	maxID += 1

	return maxID
}

func getMaxID(db *sql.DB, table string) (int, error) {
	cmd := fmt.Sprintf("SELECT MAX(id) AS max_id FROM %s", table)
	var maxID int
	err := db.QueryRow(cmd).Scan(&maxID)
	if err != nil {
		return 0, err
	}

	return maxID, nil

}
