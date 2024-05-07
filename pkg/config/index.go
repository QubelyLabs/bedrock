package config

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var cf any

func Set[T any](environment string) *T {
	if environment == "" {
		environment = gin.DebugMode
	}

	viper.AddConfigPath("env")
	viper.SetConfigName(environment)
	viper.SetConfigType("env")

	// load config from env
	viper.AutomaticEnv()

	// Viper reads all the variables from env file and log error if any found
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading env file", err)
	}

	// Viper unmarshal the loaded env varialbes into the struct
	c := new(T)
	if err := viper.Unmarshal(c); err != nil {
		log.Fatal(err)
	}

	cf = c

	return cf.(*T)
}

func Get[T any]() *T {
	return cf.(*T)
}
