package service

import (
	"ont/internal/escape"
	"os/user"
)

func Execute(table, script string) error {
	escape.LogPrint("I AM EXECUTE")
	escape.LogPrintf("Table: %s and script: %s", table, script)

	return nil

}

func ChangeNextRun() {

}

func getUID(table string) string {
	user, err := user.Lookup(table)
	if err != nil {
		escape.LogFatal(err)
	}
	userUID := user.Uid

	return userUID
}
