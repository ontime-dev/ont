/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"ont/internal/client"
	"ont/internal/dbopts"
	esc "ont/internal/escape"
	"os/user"
	"strconv"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [jobid]",
	Short: "List all the scheduled jobs",
	Long:  `Lists all the scheduled jobs for the current user.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			jobID, err := strconv.Atoi(args[0])
			if err != nil {
				esc.Error(err.Error())
			}
			err = listJobs(jobID)
			if err != nil {
				esc.Error(err.Error())
			}
		} else if len(args) == 0 {
			err := listJobs(0)
			if err != nil {
				esc.Error(err.Error())
			}
		} else {
			return errors.New("you cannot specify more than one job")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func listJobs(jobid int) error {

	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	job := []dbopts.Jobs{
		Id: jobid,
	}
	fmt.Printf("ID \t Script \t \t Next Execution Time \t Intervals \t Status \n")
	fmt.Println("----------------------------------------------------------------------------------")

	//jobs, err := dbopts.Opt("list", user.Username, job, cfgFile)
	//	fmt.Println(jobs)
	message := client.Message{
		Command: "list",
		User:    user.Username,
		Job:     job,
	}
	err = client.SendMsg(message)
	if err != nil {
		return err
	}

	//err = client.RecieveRspns()
	/*for _, job := range jobs {
		fmt.Printf("%d \t| %s \t| %s \t| %s \t \t| %s\n", job.Id, job.Script, job.Exec_time, job.Every, job.Status)
	}*/
	if err != nil {
		return err
	}
	return nil
}
