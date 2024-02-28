package dbopts

import (
	"database/sql"
	"fmt"
	"os"
)

func StopJob(db *sql.DB, user string, job Jobs) error {

	job, err := GetJob(db, user, job)

	if err != nil {
		return err
	}
	cmd := fmt.Sprintf("INSERT INTO %s (id, script, exec_time, every, status) VALUES (%d, '%s', '%s', '%s', '%s');", user, job.Id, job.Script, job.Exec_time, job.Every, job.Status)
	_, err = db.Exec(cmd)

	if err != nil {
		return err
	}

	fmt.Printf("Job %d stopped\n", job.Id)
	return nil
}

func GetJob(db *sql.DB, user string, job Jobs) (Jobs, error) {

	var status string
	cmd := fmt.Sprintf("SELECT script,exec_time,every,status FROM %s WHERE id = %d ORDER BY timestamp DESC LIMIT 1", user, job.Id)
	err := db.QueryRow(cmd).Scan(&job.Script, &job.Exec_time, &job.Every, &status)
	if err != nil {
		return job, err
	}

	if status == "Stop" {
		fmt.Println("Job is already stopped")
		os.Exit(0)
	}

	return job, nil
}
