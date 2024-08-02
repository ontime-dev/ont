/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"ont/client"
	"ont/internal/dbopts"
	"ont/internal/remote"
	"ont/internal/run"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/cobra"
)

var flags struct {
	every string
	from  string
	hour  string
	min   string
	sec   string
	day   string
	month string
	year  string
	yes   bool
}

var nodes, _ = os.Hostname()

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run <script>",
	Short: "Run the script",
	Long:  `Runs a provided script in a specific date and time and on a specific intervals that are provided as arguments.`,

	Example: `  ont run --every 1hr --from now /path/to/script.sh
  ont run --every 1d --from tomorrow /path/to/script.sh
  ont run --hour 01 /path/to/script.sh
  ont run --every 1d --from dd-MM-yyyy /path/to/script.sh
  ont run --every 1d --from dd-MM-yyyyThh:mm:ss /path/to/script.sh
  ont run --every 1d --from now -n node001 /path/to/script.sh`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := runJob(cmd, args)

		//fmt.Println(cmd.Flags().Lookup("every").Value.String())
		return err
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")
	runCmd.Flags().StringVarP(&flags.every, "every", "e", "", "Specify the period to run the script in.")
	runCmd.Flags().StringVarP(&flags.from, "from", "f", "", "From a specific time and date.")
	runCmd.Flags().StringVarP(&flags.hour, "hour", "H", "00", "Specify the hour.")
	runCmd.Flags().StringVarP(&flags.min, "min", "m", "00", "Specify the minutes.")
	runCmd.Flags().StringVarP(&flags.sec, "sec", "s", "00", "Specify the seconds.")
	runCmd.Flags().StringVarP(&flags.day, "day", "d", "", "Specify the days.")
	runCmd.Flags().StringVarP(&flags.month, "month", "M", "", "Specify the month.")
	runCmd.Flags().StringVarP(&flags.year, "year", "Y", "", "Specify the year.")
	runCmd.Flags().StringVarP(&nodes, "nodes", "n", nodes, "Specify the node list to run the job on")

	runCmd.Flags().BoolVarP(&flags.yes, "yes", "y", false, "Continue without asking for confirmation")

	runCmd.MarkFlagsRequiredTogether("every", "from")

	runCmd.MarkFlagsMutuallyExclusive("hour", "every")
	runCmd.MarkFlagsMutuallyExclusive("min", "every")
	runCmd.MarkFlagsMutuallyExclusive("sec", "every")
	runCmd.MarkFlagsMutuallyExclusive("day", "every")
	runCmd.MarkFlagsMutuallyExclusive("month", "every")
	runCmd.MarkFlagsMutuallyExclusive("year", "every")

	runCmd.MarkFlagsMutuallyExclusive("hour", "from")
	runCmd.MarkFlagsMutuallyExclusive("min", "from")
	runCmd.MarkFlagsMutuallyExclusive("sec", "from")
	runCmd.MarkFlagsMutuallyExclusive("day", "from")
	runCmd.MarkFlagsMutuallyExclusive("month", "from")
	runCmd.MarkFlagsMutuallyExclusive("year", "from")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runJob(cmd *cobra.Command, script []string) error {

	defer func() error {
		if r := recover(); r != nil {
			cmd.Usage()
			os.Exit(1)
		}

		return nil
	}()

	if len(script) != 1 {
		return errors.New("invalid number of arguments")

	}

	valid := run.CheckFromValidity(flags.from)
	if !valid {
		return errors.New("'from' time can't be in the past")
	}

	exec_time, err := run.ParseEvryFrom(flags.every, flags.from)
	if err != nil {
		return err
	}

	script_path := script[0]

	if !filepath.IsAbs(script_path) {
		script_path, err = filepath.Abs(script_path)
		if err != nil {
			return err
		}
	}

	user, err := user.Current()
	if err != nil {
		return err
	}

	// check if the script exists on the remote server
	if nodes != os.Getenv("HOSTNAME") {
		cmd := fmt.Sprintf("stat %s 2> /dev/null", script_path)
		err := remote.Run(user.Username, nodes, cmd, "", "", true, false)
		if err != nil {
			return fmt.Errorf("stat %s: no such file or directory on the remote server", script_path)
		}

	} else {
		//Check if script exists.
		_, err = os.Stat(script_path)
		if err != nil {
			return err
		}
	}

	job := dbopts.Jobs{
		Script:    script_path,
		Exec_time: exec_time,
		Every:     flags.every,
		Status:    "Active",
		RunOn:     nodes,
	}

	message := client.Message{
		Command: "run",
		User:    user.Username,
		Job:     job,
	}

	if !flags.yes {
		if flags.from != "now" {
			err := confirm(message)
			if err != nil {
				return err
			}
		}
	}

	response, err := client.SendMsg(message)
	if err != nil {
		return err
	}

	fmt.Println(response.Status)

	return nil
}

func confirm(message client.Message) error {
	var stdout_msg string

	if message.Command == "run" {
		stdout_msg = fmt.Sprintf(
			`Script: %s
Next Execution time: %s 
Interval: %s
`, message.Job.Script, message.Job.Exec_time, message.Job.Exec_time)
	} else if message.Command == "clean" {
		stdout_msg = "You are about to remove all jobs and entries.\n"
	}

	fmt.Print(stdout_msg)

	for {
		fmt.Printf("Continue?(n/Y):")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := scanner.Text()

		switch input {
		case "n", "N":
			os.Exit(0)
		case "y", "Y", "":
			return nil
		default:
			fmt.Println("Invalid Choice!")
		}
	}
}
