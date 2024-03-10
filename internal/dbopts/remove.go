package dbopts

import (
	"database/sql"
	"fmt"
)

func RemoveJob(db *sql.DB, table string, job Jobs) error {
	cmd := fmt.Sprintf("DELETE FROM %s WHERE id='%d';", table, job.Id)

	_, err := db.Exec(cmd)
	if err != nil {
		return err
	}
	fmt.Printf("Job %d removed from records.\n", job.Id)
	return nil
}
