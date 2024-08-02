package cmd

import (
	"fmt"
	"log"
	"ont/client"
	"os/user"

	"github.com/spf13/cobra"
)

var yes bool

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Removes all entries for the user running it",
	Long:  `Removes all the jobs and the table from the database for the user running it`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cleanAllJobs()

		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)

	// Here you will define your flags and configuration settings.
	cleanCmd.Flags().BoolVarP(&yes, "yes", "y", false, "Continue without asking for confirmation")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cleanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cleanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func cleanAllJobs() error {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	message := client.Message{
		Command: "clean",
		User:    user.Username,
	}

	if !yes {
		err := confirm(message)
		if err != nil {
			return err
		}
	}

	response, err := client.SendMsg(message)
	if err != nil {
		return err
	}

	fmt.Println(response.Status)

	return nil

}
