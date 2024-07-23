package dbopts

import (
	_ "github.com/go-sql-driver/mysql"
)

type Jobs struct {
	Id        int
	Script    string
	Exec_time string
	Every     string
	Status    string
	RunOn     string
}
