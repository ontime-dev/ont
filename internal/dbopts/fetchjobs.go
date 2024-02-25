package dbopts

import (
	"database/sql"
	"errors"
	"fmt"
)

func FetchJobs(db *sql.DB, user string) error {
	//WITH RankedData AS (   SELECT     golang.*,     ROW_NUMBER() OVER (PARTITION BY id ORDER BY next_run ASC) AS rn   FROM     golang ) SELECT * FROM   RankedData WHERE   rn = 1;

	// That can be deleted
	/*err := isExist(db, user)
	if err != nil {
		return err
	}*/

	cmd := fmt.Sprintf(` WITH RankedData AS (   SELECT     %s.*,     ROW_NUMBER() OVER (PARTITION BY id ORDER BY next_run ASC) 
	AS rn   FROM     %s ) SELECT * FROM   RankedData WHERE   id = 1;`, user, user)

	//fmt.Println(cmd)
	rows, err := db.Query(cmd)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {

		job := new(Jobs)
		fmt.Println(cmd)
		err = rows.Scan(&job.Id, &job.Script, &job.Next_run, &job.Every, &job.Status)
		if err != nil {
			return err
		}

		fmt.Printf("JobID: %d ----- Script: %s ----- Next Run: %s ", job.Id, job.Script, job.Next_run)

	}
	return nil

}

// THAT FUNC CAN BE DELETED //
func isExist(db *sql.DB, user string) error {
	cmd := fmt.Sprintf("SELECT * FROM %s", user)

	_, err := db.Exec(cmd)
	if err != nil {
		return errors.New("no data for that user, plesae insert a job using ont run first")
	}

	return nil

}
