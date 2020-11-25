package internal

import (
	"github.com/spf13/viper"
)

func init() {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetDefault("NumLeds", 100)
	viper.SetDefault("Brightness", 128)

	viper.ReadInConfig()
}
