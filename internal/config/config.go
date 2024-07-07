package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func LoadConfig(cfgFile string) string {
	viper.SetConfigFile(cfgFile)
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err.Error())
	}

	return viper.GetString("DBPASS")
}
