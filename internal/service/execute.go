package service

import (
	"database/sql"
	"ont/internal/dbopts"
	"ont/internal/escape"
	"ont/internal/run"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
	"time"
)

func Execute(db *sql.DB, table string, job dbopts.Jobs) {
	//Enable below when verbose
	//escape.LogPrint("Executing script", job.Script)
	escape.LogPrintf("Executing job with ID %d for user '%s'", job.Id, table)
	uid, homeDir := getUserInfo(table)
	cmd := exec.Command(job.Script)
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uid}
	cmd.Dir = homeDir
	go func() {
		err := cmd.Run()
		if err != nil {
			escape.LogPrint(err.Error())
		}
	}()
	err := ChangeExecTime(db, table, job)
	if err != nil {
		escape.LogPrint(err.Error())
	}
}

func ChangeExecTime(db *sql.DB, table string, job dbopts.Jobs) error {
	crntTime, err := time.Parse("15:04:05 Jan 02 2006", job.Exec_time)
	if err != nil {
		return err
	}

	number, last_char := run.GetLastChar(job.Every)
	job.Exec_time, err = run.ParseEvery(crntTime, last_char, number)
	if err != nil {
		return err
	}

	//Enable below when verbose
	//escape.LogPrint("Changing time to: ", job.Exec_time)
	_, err = dbopts.Insert(db, table, job, false)

	if err != nil {
		return err
	}
	return nil

}

func getUserInfo(table string) (uint32, string) {
	user, err := user.Lookup(table)
	if err != nil {
		escape.LogPrint(err)
	}
	userUID, err := strconv.ParseUint(user.Uid, 10, 32)
	if err != nil {
		escape.LogPrint(err)
	}
	return uint32(userUID), user.HomeDir
}
