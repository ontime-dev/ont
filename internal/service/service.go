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
		//logger.Fatal(err)
		escape.LogFatal(err)
	}

	defer db.Close()

	for {
		var tablename string
		var allTables []string
		tables, err := db.Query("SHOW TABLES;")
		if err != nil {
			//logger.Fatal(err)
			escape.LogFatal(err)
		}
		for tables.Next() {
			err = tables.Scan(&tablename)
			if err != nil {
				//logger.Fatal(err.Error())
				escape.LogFatal(err)
			}
			allTables = append(allTables, tablename)
		}
		var wg sync.WaitGroup
		for _, table := range allTables {
			wg.Add(1)
			go ProcessTable(db, table, &wg)
		}
		wg.Wait()
	}
}

func ProcessTable(db *sql.DB, table string, wg *sync.WaitGroup) {
	defer wg.Done()
	var job dbopts.Jobs

	maxID, err := dbopts.GetMaxID(db, table)
	if err != nil {
		//logger.Logger(err)
		escape.LogPrint(err)
	}
	//logger.Logger(table)
	escape.LogPrint(table)
	for id := 1; id <= maxID; id++ {
		job, err := dbopts.GetJob(db, table, id, job)
		//logger.Logger(job)
		escape.LogPrint(job)
		if err != nil {
			//logger.Logger(err)
			escape.LogPrint(err)
		}
		if job.Exec_time == time.Now().Format("15:04:05 Jan 02 2006") {
			//logger.Logger("Execute")
			escape.LogPrint("Execute")
			Execute(table, job.Script)
		}

	}
}
