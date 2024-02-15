package dbopts

import (
	"database/sql"
	"fmt"
)

func FetchJobs(db *sql.DB, user string) error {
	//WITH RankedData AS (   SELECT     golang.*,     ROW_NUMBER() OVER (PARTITION BY id ORDER BY next_run ASC) AS rn   FROM     golang ) SELECT * FROM   RankedData WHERE   rn = 1;

	cmd := fmt.Sprintf(`WITH RankedData AS (   SELECT     %s.*,     ROW_NUMBER() OVER (PARTITION BY id ORDER BY next_run ASC) 
	AS rn   FROM     %s ) SELECT * FROM   RankedData WHERE   rn = 1;`, user, user)

	fmt.Println(cmd)
	return nil

}
