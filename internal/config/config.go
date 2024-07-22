package config

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	CfgFile  = "/etc/ont/ont.conf"
	PassFile = "/etc/ont/ont.pass"
)

func GetConfig(parameter string) string {
	file := CfgFile

	if parameter == "DBPASS" {
		file = PassFile
	}

	viper.SetConfigFile(file)
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err.Error())
	}
	return viper.GetString(parameter)
}
