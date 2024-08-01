package dbopts

import (
	"database/sql"
	"ont/internal/escape"
)

func List(db *sql.DB, user string, verbose bool) ([]Jobs, error) {
	jobs, err := PrintJobs(db, user, verbose)
	if err != nil {
		return nil, err
	}
	if verbose {
		escape.LogPrintf("DEBUG(LST): Returning all jobs to user %s", user)
	}
	return jobs, nil

	//Can be used for showing one job. TODO.
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

func PrintJobs(db *sql.DB, table string, verbose bool) ([]Jobs, error) {
	var jobs []Jobs
	maxID, err := GetMaxID(db, table, verbose)

	if err != nil {
		return jobs, err
	}

	for id := 1; id <= maxID; id++ {
		var job Jobs
		job, err = job.GetJob(db, table, id, verbose)
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				return jobs, err
			}
		}
		if job.Status == "Active" || job.Status == "Inactive" {
			jobs = append(jobs, job)
		}
	}

	return jobs, nil
}

//NOT USED FOR THE MOMENT. TODO.
// func PrintOneJob(db *sql.DB, table string, jobid int) ([]Jobs, error) {
// 	var jobs []Jobs
// 	var script, exec_time, every, status string
// 	cmd := fmt.Sprintf("select script,exec_time,every,status from %s where id=%d ORDER BY timestamp DESC LIMIT 1;", table, jobid)
// 	err := db.QueryRow(cmd).Scan(&script, &exec_time, &every, &status)
// 	if err != nil {
// 		if err.Error() == "sql: no rows in result set" {
// 			escape.LogPrintf("Job %d doesn't exist!\n", jobid)
// 		}
// 	}
// 	job := Jobs{
// 		Id:        jobid,
// 		Script:    script,
// 		Exec_time: exec_time,
// 		Every:     every,
// 		Status:    status,
// 	}
// 	jobs = append(jobs, job)
// 	return jobs, nil

// }
