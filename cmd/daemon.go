/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/sevlyar/go-daemon"
	"github.com/spf13/cobra"
)

// daemonCmd represents the daemon command
var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Run in daemon mode",
	Long:  `Executing "ont daemon" runs the application in daemon mode`,
	Run: func(cmd *cobra.Command, args []string) {
		work()
	},
}

func init() {
	rootCmd.AddCommand(daemonCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// daemonCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// daemonCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func work() {

	cntxt := &daemon.Context{
		PidFileName: "/var/run/ont.pid",
		PidFilePerm: 0644,
		LogFileName: "/var/log/ont.log",
		LogFilePerm: 0640,
		WorkDir:     "/var/log/",
		Umask:       027,
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	log.Print("- - - - - - - - - - - - - - -")
	log.Print("daemon started")

	for {
		fmt.Println(time.Now())
		time.Sleep(2 * time.Second)

	}
}
