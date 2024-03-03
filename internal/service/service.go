package service

import (
	"database/sql"
	"fmt"
	"log"
)

func Letsgo(logger *log.Logger) error {

	db, err := sql.Open("mysql", "ont:password@/ontime")
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	cmd := fmt.Sprintf("SHOW TABLES;")
	for {
		var tablename string
		tables, err := db.Query(cmd)
		if err != nil {
			logger.Fatal(err)
		}
		go func() {
			for tables.Next() {
				err = tables.Scan(&tablename)
				if err != nil {
					logger.Fatal(err.Error())
				}
				logger.Print(tablename)
			}
		}()

	}
}
