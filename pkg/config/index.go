package config

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var cf any

func Set[T any]() *T {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	// load config from env
	viper.AutomaticEnv()

	// Get all environment variables from the server
	envVars := os.Environ()
	for _, env := range envVars {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			value := parts[1]
			viper.Set(key, value)
		}
	}

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
