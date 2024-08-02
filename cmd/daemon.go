package cmd

import (
	"github.com/ontime-dev/ont/internal/escape"
	"github.com/ontime-dev/ont/service"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// daemonCmd represents the daemon command
var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Run in daemon mode",
	Long:  `Runs the application in daemon mode`,
	Run: func(cmd *cobra.Command, args []string) {
		work()
	},
}

func init() {
	rootCmd.AddCommand(daemonCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// daemonCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// daemonCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func work() {
	escape.Init()
	escape.LogPrint(viper.ConfigFileUsed())
	//escape.LogPrint(viper.GetString("server.hostname"))

	err := service.Letsgo()
	if err != nil {
		escape.LogFatal(err.Error())
	}
}
