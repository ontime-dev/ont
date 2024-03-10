package dbopts

import (
	"database/sql"
	"errors"
	"fmt"
	"ont/internal/escape"
)

func PrintJobs(db *sql.DB, table string) error {
	maxID, err := GetMaxID(db, table)
	var job Jobs
	if err != nil {
		return err
	}
	fmt.Printf("ID \t Script \t \t Next Execution Time \t Intervals \t Status \n")
	fmt.Println("----------------------------------------------------------------------------------")

	for id := 1; id <= maxID; id++ {
		job, err = GetJob(db, table, id, job)
		if err != nil {
			return err
		}
		fmt.Printf("%d \t| %s \t| %s \t| %s \t \t| %s\n", id, job.Script, job.Exec_time, job.Every, job.Status)
	}

	return nil
}

/*
func FetchJobs(db *sql.DB, user string) error {
	//WITH RankedData AS (   SELECT     golang.*,     ROW_NUMBER() OVER (PARTITION BY id ORDER BY next_run ASC) AS rn   FROM     golang ) SELECT * FROM   RankedData WHERE   rn = 1;

	// That can be deleted
	/*err := isExist(db, user)
	if err != nil {
		return err
	}

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

}*/

// THIS FUNC CAN BE DELETED //
func isExist(db *sql.DB, table string) error {
	cmd := fmt.Sprintf("SELECT * FROM %s", table)

	_, err := db.Exec(cmd)
	if err != nil {
		return errors.New("no data for that user, plesae insert a job using ont run first")
	}

	return nil

}

func PrintOneJob(db *sql.DB, table string, jobid int) error {
	var script, exec_time, every, status string
	cmd := fmt.Sprintf("select script,exec_time,every,status from %s where id=%d ORDER BY timestamp DESC LIMIT 1;", table, jobid)
	err := db.QueryRow(cmd).Scan(&script, &exec_time, &every, &status)
	if err != nil {
		return err
	}
	fmt.Printf("ID \t Script \t \t Next Execution Time \t Intervals \t Status \n")
	fmt.Println("----------------------------------------------------------------------------------")
	fmt.Printf("%d \t| %s \t| %s \t| %s \t \t| %s\n", jobid, script, exec_time, every, status)
	escape.LogPrint("TEST TEST")

	return nil

}

func GetJob(db *sql.DB, table string, id int, job Jobs) (Jobs, error) {
	job.Id = id
	cmd := fmt.Sprintf("SELECT script,exec_time,every,status FROM %s WHERE id = %d ORDER BY timestamp DESC LIMIT 1", table, job.Id)

	err := db.QueryRow(cmd).Scan(&job.Script, &job.Exec_time, &job.Every, &job.Status)
	if err != nil {
		return job, err
	}

	return job, nil
}
