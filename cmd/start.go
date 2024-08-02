package cmd

import (
	"errors"
	"fmt"
	"ont/client"
	"ont/internal/dbopts"
	"ont/internal/run"
	"os/user"
	"strconv"

	"github.com/spf13/cobra"
)

var from string
var refresh bool

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start <jobID>",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Short: "Starts the execution of a job",
	Long:  `Marks the job with id=<jobID> as active and the job will start being executed.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("please specify the job ID")
		}

		jobid, err := strconv.Atoi(args[0])
		if err != nil {
			return errors.New("please insert a valid jobid")
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
	startCmd.Flags().BoolVarP(&refresh, "refresh", "r", false, "Refresh the job from a specific time and date. Useful when ontd service was stopped and the next exec time of a job is overdue")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func startJob(jobid int) error {
	command := "start"
	user, err := user.Current()
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
	if refresh {
		command = "refresh"
	}

	message := client.Message{
		Command: command,
		User:    user.Username,
		Job:     job,
	}

	response, err := client.SendMsg(message)
	if err != nil {
		return err
	}

	fmt.Println(response.Status)

	return nil
}
