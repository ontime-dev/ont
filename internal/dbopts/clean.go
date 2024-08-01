package dbopts

import (
	"database/sql"
	"fmt"
	"ont/internal/escape"
)

// import (
// 	"database/sql"
// 	"fmt"
// )

func CleanAllJobs(db *sql.DB, table string, verbose bool) error {
	cmd := fmt.Sprintf("Drop table %s;", table)

	if verbose {
		escape.LogPrintf("DEBUG(CLNALLJBS): cleaning all the jobs for user %s", table)
	}
	// escape.Debug(verbose, "(CLNALLJBS)", "cleaning all the jobs for user %s", table)

	_, err := db.Exec(cmd)
	if err != nil {
		return err
	}

	return nil

}
