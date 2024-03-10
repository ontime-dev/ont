package service

import (
	"database/sql"
	"ont/internal/dbopts"
	"ont/internal/escape"
	"sync"
	"time"
)

func Letsgo() error {

	db, err := sql.Open("mysql", "ont:password@/ontime")
	if err != nil {
		escape.LogFatal(err)
	}

	defer db.Close()

	for {
		var tablename string
		var allTables []string
		tables, err := db.Query("SHOW TABLES;")
		if err != nil {
			escape.LogFatal(err)
		}
		for tables.Next() {
			err = tables.Scan(&tablename)
			if err != nil {
				escape.LogFatal(err)
			}
			allTables = append(allTables, tablename)
		}
		var wg sync.WaitGroup
		for _, table := range allTables {
			wg.Add(1)
			//escape.LogPrint(table)
			go ProcessTable(db, table, &wg)
		}
		wg.Wait()
		//escape.LogPrint("I AM DONE")
		//time.Sleep(time.Second * 1)
	}
}

func ProcessTable(db *sql.DB, table string, wg *sync.WaitGroup) {
	//escape.LogPrint("Processing Table ", table)
	defer wg.Done()
	var job dbopts.Jobs

	maxID, err := dbopts.GetMaxID(db, table)
	if err != nil {
		escape.LogPrint(err)
	}
	//escape.LogPrint(table)
	for id := 1; id <= maxID; id++ {
		job, err := dbopts.GetJob(db, table, id, job)
		if err != nil {
			escape.LogPrint(err)
		}
		if job.Exec_time == time.Now().Format("15:04:05 Jan 02 2006") {
			//		escape.LogPrint("Execute")
			//CHECK IF STATUS IS ACTIVE
			if job.Status == "Active" {
				Execute(db, table, job)
				//escape.LogPrint(id)
				// err := ChangeExecTime(db, table, job)
				// if err != nil {
				// 	escape.LogFatal(err.Error())
				// }
			} else {
				escape.LogPrint("Job IS NOT ACTIVE")
			}
		}

	}
}
