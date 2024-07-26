package service

import (
	"database/sql"
	"fmt"
	"ont/internal/config"
	"ont/internal/dbopts"
	"ont/internal/escape"
	"sync"
	"time"
)

func Letsgo() error {
	password, ip, port := config.GetConfig("DBPASS"), config.GetConfig("SERVER_IP"), config.GetConfig("SERVER_PORT")
	pass_cmd := fmt.Sprintf("ont:%s@/ontime", password)

	db, err := sql.Open("mysql", pass_cmd)
	if err != nil {
		escape.LogFatal(err)
	}

	go Server(db, ip, port)
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
		//verbose logging
		//escape.LogPrint("(PrcsTbl):", err)
		escape.LogPrint(err)
	}

	for id := 1; id <= maxID; id++ {
		job, err := job.GetJob(db, table, id)
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				escape.LogPrint(err)
			}

		}
		if job.Exec_time == time.Now().Format("15:04:05 Jan 02 2006") {
			Execute(db, table, job)
			//CHECK IF STATUS IS ACTIVE
			// if job.Status == "Active" {
			// 	Execute(db, table, job)
			// } else {
			// 	escape.LogPrint("Job IS NOT ACTIVE")
			// }
		}

	}
}
