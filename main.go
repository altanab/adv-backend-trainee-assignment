package main

import (
	"adv/internal/app"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var ConfigPath = "config.json"

func main() {
	var config app.Config
	configFile, err := os.Open(ConfigPath)
	if err != nil {
		log.Fatalf("cannot open config file: %s", err.Error())
	}

	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		log.Fatalf("cannot unmarshal config file: %s", err.Error())
	}
	configFile.Close()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	app.StartServer(quit, config)
}

