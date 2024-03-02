/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"ont/internal/dbopts"
	"os/user"
	"strconv"

	"github.com/spf13/cobra"
)

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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func startJob(jobid int) error {
	user, err := user.Current()
	if err != nil {
		return err
	}

	job := dbopts.Jobs{
		Id:     jobid,
		Status: "Active",
	}
	err = dbopts.Opt("start", user.Username, job)
	if err != nil {
		return err
	}

	return nil
}
