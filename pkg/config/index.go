package config

import (
	"log"

	"github.com/spf13/viper"
)

var cf any

func Set[T any]() *T {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	// load config from env
	viper.AutomaticEnv()

	// Viper reads all the variables from env file and log error if any found
	if err := viper.ReadInConfig(); err != nil {
		log.Println("Error reading env file", err)
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
