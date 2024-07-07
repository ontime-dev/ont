/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	esc "ont/internal/escape"
	"strconv"

	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove [jobid]",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			jobID, err := strconv.Atoi(args[0])
			if err != nil {
				esc.Error(err.Error())
			}
			err = removeJob(jobID)
			if err != nil {
				esc.Error(err.Error())
			}
		} else if len(args) == 0 {
			return errors.New("please specify one job to remove")
		} else {
			return errors.New("you cannot remove more than one job at a time")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func removeJob(jobid int) error {
	/*user, err := user.Current()
	if err != nil {
		return err
	}
	job := dbopts.Jobs{
		Id: jobid,
	}

	//err = dbopts.Opt("remove", user.Username, job, cfgFile)
	dbopts.Opt("remove", user.Username, job, cfgFile)

	if err != nil {
		return err
	}*/
	return nil
}
