/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"

	"ont/internal/service"

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
	/*
		cntxt := &daemon.Context{
			PidFileName: "/var/run/ont.pid",
			PidFilePerm: 0644,
			LogFileName: "/var/log/ont.log",
			LogFilePerm: 0640,
			WorkDir:     "/",
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
	*/
	logFile, err := os.OpenFile("/var/log/ont.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		log.Fatal("Error")
	}
	defer logFile.Close()

	logger := log.New(logFile, "PREFIX:", log.Ldate|log.Ltime|log.Lshortfile)
	err = service.Letsgo(logger)

	if err != nil {
		logger.Fatal(err.Error())
	}
}
