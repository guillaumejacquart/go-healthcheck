package pkg

import (
	"github.com/spf13/viper"
)

func initConfig(configPath string) {
	viper.SetDefault("port", 8080)
	viper.SetDefault("history.enabled", "false")
	viper.SetDefault("db.type", "sqlite3")
	viper.SetDefault("db.path", "data.db")

	viper.SetConfigName("config")                // name of config file (without extension)
	viper.AddConfigPath("/etc/go-healthcheck/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.go-healthcheck") // call multiple times to add many search paths
	viper.AddConfigPath(configPath)              // optionally look for config in the working directory
	viper.AutomaticEnv()

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(err)
	}
}

func Run(configPath string) {
	initConfig(configPath)

	initDb()
	go registerChecks()

	server := createServer()
	server.initializeMiddlewares()
	server.setupRoutes()
	server.serve(viper.GetInt("port"))
}
