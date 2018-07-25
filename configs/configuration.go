// Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

package configs

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

// Configuration is the top level configuration data from
// the config file.
type Configuration struct {
	Backend        BackendConfiguration
	Authentication AuthenticationConfiguration
	Client         ClientConfiguration
}

// LoadConfig loads the config file into the configuration
// objects and returns the top level configuration.
func LoadConfig() (configuration Configuration) {
	// TODO: What about pflags?
	// TODO: What about env vars?
	// viper.SetConfigName("config") // TODO: ENV_VAR to specify which config file to use?
	viper.SetConfigName("config-inttest")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.WithError(err).Fatal("Error reading config file.")
		os.Exit(-1)
	}
	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.WithError(err).Fatal("Unable to decode into struct.")
		os.Exit(-1)
	}

	// TODO: Change client creds to load from viper
	// var c ClientConfiguration
	// c.
	// file, err := ioutil.ReadFile("client.json")
	// if err != nil {
	// 	log.WithError(err).Fatal("Failed to load client configuration.")
	// 	os.Exit(-1)
	// }
	// json.Unmarshal(file, &c)
	configuration.Client = ClientConfiguration{
		ID:     "id",
		Secret: "secret",
	}
	return
}
