package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var cfgFile = "/etc/ont/ont.conf"

func LoadConfig() (password, port string) {
	viper.SetConfigFile(cfgFile)
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err.Error())
	}

	return viper.GetString("DBPASS"), viper.GetString("SERVER_PORT")
}
