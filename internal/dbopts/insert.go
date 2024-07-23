package dbopts

import (
	"database/sql"
	"fmt"
)

func Insert(db *sql.DB, user string, job Jobs, new bool) (int, error) {

	var id int
	err := Create(db, user)
	if err != nil {
		return 0, err
	}

	if new {
		id = setID(db, user)
	} else {
		id = job.Id
		//Enable below when verbose
		//escape.LogPrint("Inserting in table")
	}

	cmd := fmt.Sprintf("INSERT INTO %s (id, script, exec_time, every, status) VALUES (%d, '%s', '%s', '%s', '%s');", user, id, job.Script, job.Exec_time, job.Every, job.Status)

	_, err = db.Exec(cmd)
	if err != nil {
		return 0, err
	}

	return id, err

}

func setID(db *sql.DB, table string) int {

	maxID, err := GetMaxID(db, table)
	if err != nil {
		if maxID == 0 {
			return 1
		}
	}
	maxID += 1

	return maxID
}

func GetMaxID(db *sql.DB, table string) (int, error) {
	cmd := fmt.Sprintf("SELECT MAX(id) AS max_id FROM %s", table)
	var maxID int
	err := db.QueryRow(cmd).Scan(&maxID)
	if err != nil {
		return 0, err
	}

	return maxID, nil

}
