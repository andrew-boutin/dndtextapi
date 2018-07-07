package configs

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Configuration struct {
	Backend BackendConfiguration
}

func LoadConfig() (configuration Configuration) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
		os.Exit(-1)
	}
	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %s", err)
		os.Exit(-1)
	}
	return
}
