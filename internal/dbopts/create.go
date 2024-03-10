package dbopts

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const table = `(
	timestamp timestamp(3),
	id int,
	script varchar(255),
	exec_time varchar(255),
	every varchar(255),
	status varchar(255)
)`

func Create(db *sql.DB, user string) error {

	cmd := "CREATE TABLE IF NOT EXISTS " + user + table + ";"
	//fmt.Println(cmd)

	_, err := db.Exec(cmd)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil

}

/*
func IsExist(err error) {
	if
}*/
