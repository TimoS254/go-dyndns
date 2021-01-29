package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	IntervalMinutes int8
	Domains         []Domain
}

func generateConfig() {
	file, _ := os.Create("config.json")
	defer file.Close()
	data, _ := json.MarshalIndent(config, "", "  ")
	file.Write([]byte(data))
}
