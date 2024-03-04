package service

import (
	"database/sql"
	"log"
	"ont/internal/dbopts"
	"sync"
	"time"
)

func Letsgo(logger *log.Logger) error {

	db, err := sql.Open("mysql", "ont:password@/ontime")
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()

	for {
		var tablename string
		var allTables []string
		tables, err := db.Query("SHOW TABLES;")
		if err != nil {
			logger.Fatal(err)
		}
		for tables.Next() {
			err = tables.Scan(&tablename)
			if err != nil {
				logger.Fatal(err.Error())
			}
			allTables = append(allTables, tablename)
		}
		var wg sync.WaitGroup
		for _, table := range allTables {
			wg.Add(1)
			go ProcessTable(db, table, logger, &wg)
		}
		wg.Wait()
	}
}

func ProcessTable(db *sql.DB, table string, logger *log.Logger, wg *sync.WaitGroup) {
	defer wg.Done()
	var job dbopts.Jobs

	maxID, err := dbopts.GetMaxID(db, table)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println(table)
	for id := 1; id <= maxID; id++ {
		job, err := dbopts.GetJob(db, table, id, job)
		logger.Print(job.Id)
		if err != nil {
			logger.Fatal(err)
		}
		if job.Exec_time == time.Now().Format("15:04:05 Jan 02 2006") {
			logger.Println("Execute")
		}

	}
}
