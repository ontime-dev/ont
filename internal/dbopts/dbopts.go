package dbopts

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/ontime-dev/ont/internal/escape"

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

func (job Jobs) Insert(db *sql.DB, user string, new, verbose bool) (int, error) {

	err := Create(db, user, verbose)
	if err != nil {
		return 0, err
	}

	var id int
	var logMessage string
	if new {
		id = setID(db, user, verbose)
		logMessage = fmt.Sprintf("INFO: New entry inserted for job %d", id)
	} else {
		id = job.Id
		logMessage = fmt.Sprintf("INFO: Job %d updated for user %s", id, user)
	}

	cmd := fmt.Sprintf("INSERT INTO %s (id, script, exec_time, every, status, runon) VALUES (%d, '%s', '%s', '%s', '%s', '%s');", user, id, job.Script, job.Exec_time, job.Every, job.Status, job.RunOn)
	// Enable below when verbose
	if verbose {
		escape.LogPrint("DEBUG(INSRT): Inserting new entry in table")
	}
	_, err = db.Exec(cmd)
	if err != nil {
		return 0, err
	}

	escape.LogPrintf(logMessage)

	return id, err

}

func (job Jobs) RemoveJob(db *sql.DB, table string, verbose bool) error {
	cmd := fmt.Sprintf("DELETE FROM %s WHERE id='%d';", table, job.Id)

	job, err := job.GetJob(db, table, job.Id, verbose)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			err_msg := fmt.Sprintf("Job %d doesn't exist.", job.Id)
			return errors.New(err_msg)
		}
	}
	if verbose {
		escape.LogPrintf("DEBUG(RMVJB): Removing Job %d for user %s", job.Id, table)
	}
	_, err = db.Exec(cmd)
	if err != nil {
		return err
	}

	//For Verbose logging
	if verbose {
		escape.LogPrintf("DEBUG(RMVJB): Job %d removed from records for user %s.\n", job.Id, table)
	}
	return nil
}

func (job Jobs) GetJob(db *sql.DB, table string, id int, verbose bool) (Jobs, error) {
	job.Id = id
	cmd := fmt.Sprintf("SELECT script,exec_time,every,status,runon FROM %s WHERE id = %d ORDER BY timestamp DESC LIMIT 1;", table, job.Id)

	err := db.QueryRow(cmd).Scan(&job.Script, &job.Exec_time, &job.Every, &job.Status, &job.RunOn)
	if err != nil {
		return job, err
	}

	return job, nil
}

func (job Jobs) ChangeJobStatus(db *sql.DB, user, Jobstatus string, refresh, verbose bool) error {

	job, err := job.GetJobStatus(db, user, Jobstatus, refresh, verbose)
	if err != nil {
		return err
	}
	if verbose {
		escape.LogPrintf("DEBUG(CHNGJBSTTS): Changing the status of job %d for user %s.", job.Id, user)
	}
	_, err = job.Insert(db, user, false, verbose)

	if err != nil {
		return err
	}
	return nil
}

func (job Jobs) GetJobStatus(db *sql.DB, table, Jobstatus string, refresh, verbose bool) (Jobs, error) {

	if verbose {
		escape.LogPrintf("DEBUG(GTJBSTTS): Retrieving the current status of the job %d", job.Id)
	}
	oldjob, err := job.GetJob(db, table, job.Id, verbose)
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

func setID(db *sql.DB, table string, verbose bool) int {
	maxID, err := GetMaxID(db, table, verbose)
	if err != nil {
		if maxID == 0 {
			return 1
		}
	}
	maxID += 1

	if verbose {
		escape.LogPrintf("DEBUG(STID): Assigning new ID=%d to new job for user %s.", maxID, table)
	}
	return maxID
}

func GetMaxID(db *sql.DB, table string, verbose bool) (int, error) {

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
