/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
)

var from string

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start <jobID>",
	Short: "Starts the execution of a job",
	Long:  `Running 'ont start' will mark the job <jobID> as active and the job will start being executed`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("please specify the job ID")
		}

		jobid, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		err = startJob(jobid)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().StringVarP(&from, "from", "f", "now", "Start the job from a specific time and date.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func startJob(jobid int) error {
	/*user, err := user.Current()
	if err != nil {
		return err
	}

	crnttime, err := run.ParseFrom(from)

	if err != nil {
		return err
	}
	exec_time := crnttime.Format("15:04:05 Jan 02 2006")
	job := dbopts.Jobs{
		Id:        jobid,
		Status:    "Active",
		Exec_time: exec_time,
	}
	dbopts.Opt("start", user.Username, job, cfgFile)
	if err != nil {
		return err
	}
	*/
	return nil
}
