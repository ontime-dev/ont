package service

import (
	"database/sql"
	"ont/internal/dbopts"
	"ont/internal/escape"
	"ont/internal/remote"
	"ont/internal/run"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
	"time"
)

func Execute(db *sql.DB, table string, job dbopts.Jobs, verbose bool) {
	//Enable below when verbose
	//escape.LogPrint("Executing script", job.Script)
	escape.LogPrintf("Executing job with ID %d for user '%s' on %s", job.Id, table, job.RunOn)
	if verbose {
		escape.LogPrintf("DEBUG(EXCT): Getting user '%s' information", table)
	}
	uid, homeDir := getUserInfo(table)
	cmd := exec.Command(job.Script)
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uid}
	cmd.Dir = homeDir
	go func() {
		if job.RunOn == os.Getenv("HOSTNAME") {
			if verbose {
				escape.LogPrint("DEBUG(EXCT): Executing the job Locally")
			}
			err := cmd.Run()
			if err != nil {
				escape.LogPrint(err.Error())
			}
		} else {
			if verbose {
				escape.LogPrintf("DEBUG(EXCT): Executing the job on the remote server %s", job.RunOn)
			}
			err := remote.Run(table, job.RunOn, job.Script, "", "", false, verbose)
			if err != nil {
				escape.LogPrint(err.Error())
			}
		}
	}()
	err := ChangeExecTime(db, table, job, verbose)
	if err != nil {
		escape.LogPrint(err.Error())
	}
}

func ChangeExecTime(db *sql.DB, table string, job dbopts.Jobs, verbose bool) error {
	crntTime, err := time.Parse("15:04:05 Jan 02 2006", job.Exec_time)
	if err != nil {
		return err
	}
	if verbose {
		escape.LogPrintf("DEBUG(CHNGEXCTIM): Changing the next executing time for job %d for the user %s", job.Id, table)
		escape.LogPrintf("DEBUG(CHNGEXCTIM): Parsing 'every' flag")
	}

	number, last_char := run.GetLastChar(job.Every)
	job.Exec_time, err = run.ParseEvery(crntTime, last_char, number)
	if err != nil {
		return err
	}

	//Enable below when verbose
	//escape.LogPrint("Changing time to: ", job.Exec_time)

	if verbose {
		escape.LogPrintf("DEBUG(CHNGEXCTIM): Changed next execution time for job %d for the user %s to %s", job.Id, table, job.Exec_time)
	}
	_, err = job.Insert(db, table, false, verbose)

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
