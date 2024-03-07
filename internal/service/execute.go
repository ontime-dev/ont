package service

import (
	"ont/internal/escape"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
)

func Execute(table, script string) error {

	uid, homeDir := getUserInfo(table)
	escape.LogPrint(script)
	cmd := exec.Command(script)
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uid}
	cmd.Dir = homeDir
	err := cmd.Run()
	if err != nil {
		escape.LogPrint(err.Error())
	}
	return nil

}

func ChangeNextRun() {

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
