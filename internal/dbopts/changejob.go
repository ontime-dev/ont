package dbopts

import (
	"database/sql"
	"fmt"
	"ont/internal/escape"
	"strings"
)

func ChangeJobStatus(db *sql.DB, user, Jobstatus string, job Jobs) error {

	job, err := GetJobStatus(db, user, Jobstatus, job)
	if err != nil {
		return err
	}

	cmd := fmt.Sprintf("INSERT INTO %s (id, script, exec_time, every, status) VALUES (%d, '%s', '%s', '%s', '%s');", user, job.Id, job.Script, job.Exec_time, job.Every, job.Status)
	_, err = db.Exec(cmd)

	if err != nil {
		return err
	}

	/*	var status string
		if Jobstatus == "Inactive" {
			status = "stopped"
		} else {
			status = "started"
		}*/
	Jobstatus = strings.ToLower(Jobstatus[:1]) + Jobstatus[1:]
	fmt.Printf("Job %d is now %s.\n", job.Id, Jobstatus)
	return nil
}

func GetJobStatus(db *sql.DB, table, Jobstatus string, job Jobs) (Jobs, error) {
	/*
		var status string

		cmd := fmt.Sprintf("SELECT script,exec_time,every,status FROM %s WHERE id = %d ORDER BY timestamp DESC LIMIT 1", table, job.Id)
		err := db.QueryRow(cmd).Scan(&job.Script, &job.Exec_time, &job.Every, &status)
		if err != nil {
			return job, err
		}*/

	oldjob, err := GetJob(db, table, job.Id, job)
	if err != nil {
		return job, err
	}

	if oldjob.Status == Jobstatus {
		Jobstatus = strings.ToLower(Jobstatus[:1]) + Jobstatus[1:]
		err := fmt.Sprintf("Job %d is already %s.", job.Id, Jobstatus)
		escape.Error(err)
	}

	oldjob.Status = job.Status

	return oldjob, nil
}
