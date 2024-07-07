/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

type OntCfg interface {
	GetCfg() string
}

//const customUsage = `Usage:
//ont [command] <script>`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ont",
	Args:  cobra.MatchAll(cobra.ExactArgs(2), cobra.OnlyValidArgs),
	Short: "Schedule your jobs and run them",
	Long:  `ont is a CLI application to schedule running a script or a job on the node. It is an extended version of cronjobs.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	//cobra.OnInitialize(initConfig)
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	//rootCmd.SetUsageFunc(Usage)
	//rootCmd.SetUsageTemplate(customUsage)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "/etc/ont/ont.conf", "config file")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}

// initConfig reads in config file and ENV variables if set.
// func initConfig() {

// 	viper.SetConfigFile(cfgFile)
// 	viper.SetConfigType("env")

// 	err := viper.ReadInConfig()
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}

// }

func GetCfg() string {
	return cfgFile
}
