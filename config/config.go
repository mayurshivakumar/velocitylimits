package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Configurations struct {
	VelocityLimit VelocityLimit
}

type VelocityLimit struct {
	MaxDailyLoadLimit    float64
	MaxDailyTransactions int
	MaxWeeklyLoadLimit   float64
	InputFile            string
	OutputFile           string
}

// TODO: test me please!!!
//ParseConfig ...
func ParseConfig() *Configurations {
	var config Configurations
	viper.SetConfigName("config")
	viper.AddConfigPath("../config")
	viper.SetConfigType("yml")
	// TODO: Path for the file needs to be handled better
	if err := viper.ReadInConfig(); err != nil {
		logrus.Panicf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&config)
	if err != nil {
		logrus.Panicf("Unable to decode into struct, %v", err)
	}

	return &config
}
