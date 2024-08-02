package cmd

import (
	"errors"
	"fmt"
	"os/user"
	"strconv"

	"github.com/ontime-dev/ont/client"
	"github.com/ontime-dev/ont/internal/dbopts"

	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop <jobID>",
	Short: "Stops the execution of a job",
	Long:  `Marks the job with id=<jobID> as inactive and the job will be suspended.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("please specify the job ID")
		}

		jobid, err := strconv.Atoi(args[0])
		if err != nil {
			return errors.New("please insert a valid jobid")
		}

		err = stopJob(jobid)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stopCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stopCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func stopJob(jobid int) error {
	user, err := user.Current()
	if err != nil {
		return err
	}

	job := dbopts.Jobs{
		Id:        jobid,
		Status:    "Inactive",
		Exec_time: "Not Available",
	}
	message := client.Message{
		Command: "stop",
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
