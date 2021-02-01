package main

import (
	"UserTransferService/src/db_service"
	"UserTransferService/src/system/config"
	"UserTransferService/src/system/console"
	"UserTransferService/src/system/l2f"
	"UserTransferService/src/system/router"
	"os"
)

const configFile = "config.json"

func init() {
	// Initializes logger file
	l2f.Log = l2f.InitLogger()
	// Reads config
	if err := config.ReadConfig(configFile); err != nil {
		l2f.Log.Println(err)
		return
	} else {
		l2f.Log.Println("Successfully loaded a config")
	}
	if err := db_service.Connect(); err != nil {
		l2f.Log.Println(err)
		return
	} else {
		l2f.Log.Println("Successfully connected to db")
	}
}

func main() {
	if dockerEnv := os.Getenv("DOCKER"); dockerEnv != "true" {
		// Starts reading console
		go console.ReadConsole()
	} else {
		config.CF.DB.Domain = "postgre_bal"
		if err := db_service.Connect(); err != nil {
			l2f.Log.Println(err)
			return
		} else {
			l2f.Log.Println("Successfully connected to db")
		}
		// Starts listening http requests
		router.LaunchRouter()
	}
}