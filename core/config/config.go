package config

import (
	"log"

	"github.com/spf13/viper"
)

var (
	SessionName  string
	LoginVersion string
)

func GetString(key string) string {
	return viper.GetString(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func InitConfig(configFile string) {
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Open config file error - %v, Please use -c <configFile> and ensure the config file is correct", err)
	}
	SessionName = viper.GetString("session.name")
	if SessionName == "" {
		SessionName = "_dpm_sessionid_"
	}
	LoginVersion = viper.GetString("login.version")
	if LoginVersion == "" {
		LoginVersion = "dpm"
	}
}
