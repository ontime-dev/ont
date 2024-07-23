package dbopts

import (
	"database/sql"
	"fmt"
	"ont/internal/escape"
)

func List(db *sql.DB, user string) ([]Jobs, error) {
	jobs, err := PrintJobs(db, user)
	if err != nil {
		return nil, err
	}
	return jobs, nil

	//Can be used for showing one job.
	/*
		if job.Id == 0 {
			err, jobs := PrintJobs(db, user)
			if err != nil {
				return err, nil
			}
			return nil, jobs
		} /*else if job.Id > 0 {
			err, jobs := PrintOneJob(db, user, job.Id)
			if err != nil {
				return err, nil
			}
			return nil, jobs
		} else {
			return errors.New("please enter a valid jobID"), nil
		}*/
}

func PrintJobs(db *sql.DB, table string) ([]Jobs, error) {
	var jobs []Jobs
	maxID, err := GetMaxID(db, table)
	var job Jobs
	if err != nil {
		return jobs, err
	}
	escape.LogPrint(maxID)

	for id := 1; id <= maxID; id++ {

		job, err = GetJob(db, table, id, job)
		escape.LogPrint(job)
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				return jobs, err
			}
		}
		jobs = append(jobs, job)

	}

	return jobs, nil
}

func PrintOneJob(db *sql.DB, table string, jobid int) ([]Jobs, error) {
	var jobs []Jobs
	var script, exec_time, every, status string
	cmd := fmt.Sprintf("select script,exec_time,every,status from %s where id=%d ORDER BY timestamp DESC LIMIT 1;", table, jobid)
	err := db.QueryRow(cmd).Scan(&script, &exec_time, &every, &status)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			escape.LogPrintf("Job %d doesn't exist!\n", jobid)
		}
	}
	job := Jobs{
		Id:        jobid,
		Script:    script,
		Exec_time: exec_time,
		Every:     every,
		Status:    status,
	}
	jobs = append(jobs, job)
	return jobs, nil

}

func GetJob(db *sql.DB, table string, id int, job Jobs) (Jobs, error) {
	job.Id = id
	cmd := fmt.Sprintf("SELECT script,exec_time,every,status,runon FROM %s WHERE id = %d ORDER BY timestamp DESC LIMIT 1;", table, job.Id)

	err := db.QueryRow(cmd).Scan(&job.Script, &job.Exec_time, &job.Every, &job.Status, &job.RunOn)
	if err != nil {
		return job, err
	}

	return job, nil
}
