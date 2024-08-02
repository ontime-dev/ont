package service

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ontime-dev/ont/internal/config"
	"github.com/ontime-dev/ont/internal/dbopts"
	"github.com/ontime-dev/ont/internal/escape"
)

func Letsgo() error {
	var verbose bool

	if strings.ToLower(config.GetConfig("DEBUG")) == "true" {
		verbose = true
	} else if strings.ToLower(config.GetConfig("DEBUG")) == "false" {
		verbose = false
	} else {
		escape.LogFatal("Invalid value of DEBUG (True/False)")
	}

	password, ip, port := config.GetConfig("DBPASS"), config.GetConfig("SERVER_IP"), config.GetConfig("SERVER_PORT")
	pass_cmd := fmt.Sprintf("ont:%s@/ontime", password)

	db, err := sql.Open("mysql", pass_cmd)
	if err != nil {
		escape.LogFatal(err)
	}

	go Server(db, ip, port, verbose)
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

			go ProcessTable(db, table, &wg, verbose)
		}
		wg.Wait()
	}
}

func ProcessTable(db *sql.DB, table string, wg *sync.WaitGroup, verbose bool) {

	defer wg.Done()
	var job dbopts.Jobs

	maxID, err := dbopts.GetMaxID(db, table, verbose)
	if err != nil {
		if err.Error() == fmt.Sprintf("Table 'ontime.%s' doesn't exist", table) {
			return
		}
		//verbose logging
		//escape.LogPrint("(PrcsTbl):", err)
		escape.LogPrint(err)
	}

	for id := 1; id <= maxID; id++ {
		job, err := job.GetJob(db, table, id, verbose)
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				escape.LogPrint(err)
			}

		}
		if job.Exec_time == time.Now().Format("15:04:05 Jan 02 2006") {
			Execute(db, table, job, verbose)
		}

	}

}
