/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
)

//	var refreshFrom string

// restartCmd represents the restart command
var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("please specify the job ID")
		}

		jobid, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		err = startJob(jobid, true)
		if err != nil {
			return err
		}
		return nil

	},
}

func init() {
	rootCmd.AddCommand(refreshCmd)

	//startCmd.Flags().StringVarP(&refreshFrom, "froms", "f", "now", "Start the job from a specific time and date.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// restartCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// restartCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// func restartJob(jobid int) error {
// 	user, err := user.Current()
// 	if err != nil {
// 		return err
// 	}
// 	crnttime, err := run.ParseFrom(from)

// 	if err != nil {
// 		return err
// 	}
// 	exec_time := crnttime.Format("15:04:05 Jan 02 2006")

// 	job := dbopts.Jobs{
// 		Id:        jobid,
// 		Status:    "Active",
// 		Exec_time: exec_time,
// 	}
// 	message := client.Message{
// 		Command: "restart",
// 		User:    user.Username,
// 		Job:     job,
// 	}

// 	response, err := client.SendMsg(message)
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Println(response.Status)

// 	return nil
// }
