/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"

	"ont/internal/dbopts"
	"ont/internal/run"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/cobra"
)

// var evFrm, timedate bool
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

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run <script>",
	Short: "Run the script",
	Long: `This command runs the script in the specified time. 

For example:
  ont run --every 1hr --from now /path/to/script.sh
  ont run --every 1d --from tomorrow /path/to/script.sh 
  ont run --hour 01 /path/to/script.sh
  ont run --every 1d --from dd-MM-yyyy /path/to/script.sh
  ont run --every 1d --from dd-MM-yyyyThh:mm:ss /path/to/script.sh`,
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

	runCmd.Flags().BoolVarP(&flags.yes, "yes", "y", false, "Continue with asking for confirmation")

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

	if len(script) != 1 {
		return errors.New("invalid number of arguments")

	}

	exec_time, err := run.ParseEvryFrom(flags.every, flags.from)
	if err != nil {
		return err
	}

	script_path := script[0]

	if !filepath.IsAbs(script_path) {
		script_path, _ = filepath.Abs(script_path)
		fmt.Println("script: ", script_path)
	}

	//Check if script exists.
	_, err = os.Stat(script_path)
	if err != nil {
		return err
	}

	//wd, err := os.Getwd()
	if err != nil {
		return err
	}

	//fullpath := os.Get
	/*script := exec.Command(filepath[0])

	err = script.Run()
	if err != nil {
		return err
	}*/

	user, err := user.Current()
	if err != nil {
		return err
	}

	//err = dbopts.Create(user.Username)
	/*
		job := new(Jobs)
		job.Script = script_path
		job.Next_run = next_run
	*/
	job := dbopts.Jobs{
		Script:    script_path,
		Exec_time: exec_time,
		Every:     flags.every,
	}

	if !flags.yes {
		confirm(job)
	}

	err = dbopts.Opt("insert", user.Username, job)
	if err != nil {
		return err
	}

	return nil
}

func confirm(job dbopts.Jobs) {
	fmt.Printf(
		`Script: %s
Next Execution time: %s 
Interval: %s 
Continue?(n/Y):`, job.Script, job.Exec_time, job.Every)
}
