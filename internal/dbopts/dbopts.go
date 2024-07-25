package dbopts

import (
	"database/sql"
	"errors"
	"fmt"
	"ont/internal/escape"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Jobs struct {
	Id        int
	Script    string
	Exec_time string
	Every     string
	Status    string
	RunOn     string
}

func (job Jobs) Insert(db *sql.DB, user string, new bool) (int, error) {

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

	cmd := fmt.Sprintf("INSERT INTO %s (id, script, exec_time, every, status, runon) VALUES (%d, '%s', '%s', '%s', '%s', '%s');", user, id, job.Script, job.Exec_time, job.Every, job.Status, job.RunOn)
	_, err = db.Exec(cmd)
	if err != nil {
		return 0, err
	}

	escape.LogPrintf("New entry inserted for job %d", id)

	return id, err

}

func (job Jobs) RemoveJob(db *sql.DB, table string) error {
	cmd := fmt.Sprintf("DELETE FROM %s WHERE id='%d';", table, job.Id)

	job, err := job.GetJob(db, table, job.Id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			err_msg := fmt.Sprintf("Job %d doesn't exist.", job.Id)
			return errors.New(err_msg)
		}
	}
	_, err = db.Exec(cmd)
	if err != nil {
		return err
	}

	//For Verbose logging
	//escape.LogPrint("Job %d removed from records.\n", job.Id)

	return nil
}

func (job Jobs) GetJob(db *sql.DB, table string, id int) (Jobs, error) {
	job.Id = id
	cmd := fmt.Sprintf("SELECT script,exec_time,every,status,runon FROM %s WHERE id = %d ORDER BY timestamp DESC LIMIT 1;", table, job.Id)

	err := db.QueryRow(cmd).Scan(&job.Script, &job.Exec_time, &job.Every, &job.Status, &job.RunOn)
	if err != nil {
		return job, err
	}

	return job, nil
}

func (job Jobs) ChangeJobStatus(db *sql.DB, user, Jobstatus string, refresh bool) error {

	job, err := job.GetJobStatus(db, user, Jobstatus, refresh)
	if err != nil {
		return err
	}

	cmd := fmt.Sprintf("INSERT INTO %s (id, script, exec_time, every, status) VALUES (%d, '%s', '%s', '%s', '%s');", user, job.Id, job.Script, job.Exec_time, job.Every, job.Status)
	_, err = db.Exec(cmd)

	if err != nil {
		return err
	}
	return nil
}

func (job Jobs) GetJobStatus(db *sql.DB, table, Jobstatus string, refresh bool) (Jobs, error) {

	oldjob, err := job.GetJob(db, table, job.Id)
	if err != nil {
		return job, err
	}

	//The below allows the users to see message "job is already active/inactive"

	if !refresh {
		if oldjob.Status == Jobstatus {
			Jobstatus = strings.ToLower(Jobstatus[:1]) + Jobstatus[1:]
			err := fmt.Sprintf("Job %d is already %s.", job.Id, Jobstatus)
			return job, errors.New(err)
		}
	}

	oldjob.Status = job.Status
	oldjob.Exec_time = job.Exec_time

	return oldjob, nil
}

func setID(db *sql.DB, table string) int {

	maxID, err := GetMaxID(db, table)
	if err != nil {
		if maxID == 0 {
			return 1
		}
	}
	maxID += 1

	//verbose logging
	escape.LogPrintf("Job Id=%d assigned to new job.", maxID)

	return maxID
}

func GetMaxID(db *sql.DB, table string) (int, error) {
	cmd := fmt.Sprintf("SELECT MAX(id) AS max_id FROM %s", table)
	maxID := 0
	err := db.QueryRow(cmd).Scan(&maxID)
	if err != nil {
		if strings.Contains(err.Error(), "converting NULL to int is unsupported") {
			return 0, nil
		}
		return 0, err
	}

	return maxID, nil

}
