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
	escape.LogPrint("Executing script", job.Script)
	uid, homeDir := getUserInfo(table)
	escape.LogPrint(job.Script)
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
	//escape.LogPrint(job)
	//escape.LogPrint(job.Exec_time)
	crntTime, err := time.Parse("15:04:05 Jan 02 2006", job.Exec_time)
	if err != nil {
		return err
	}
	//escape.LogPrint(crntTime.Format("15:04:05 Jan 02 2006"))
	number, last_char := run.GetLastChar(job.Every)
	job.Exec_time, err = run.ParseEvery(crntTime, last_char, number)
	if err != nil {
		return err
	}

	escape.LogPrint("Changing time to: ", job.Exec_time)
	err = dbopts.Insert(db, table, job, false)

	if err != nil {
		return err
	}
	escape.LogPrint(job.Exec_time)
	return nil

}

func getUserInfo(table string) (uint32, string) {
	user, err := user.Lookup(table)
	if err != nil {
		escape.LogFatal(err)
	}
	userUID, err := strconv.ParseUint(user.Uid, 10, 32)
	if err != nil {
		escape.LogFatal(err)
	}
	return uint32(userUID), user.HomeDir
}
