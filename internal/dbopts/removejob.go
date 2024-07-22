package dbopts

import (
	"database/sql"
	"errors"
	"fmt"
)

func RemoveJob(db *sql.DB, table string, job Jobs) error {
	cmd := fmt.Sprintf("DELETE FROM %s WHERE id='%d';", table, job.Id)

	job, err := GetJob(db, table, job.Id, job)
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
