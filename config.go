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
	file, _ := os.Create("config-default.json")
	defer file.Close()
	exampleConfig := Config{
		IntervalMinutes: 5,
		Domains: []Domain{{
			DomainName:     "example.com",
			IP4:            true,
			IP6:            true,
			APIToken:       "yourAPIToken",
			ZoneIdentifier: "yourZoneIdentifier",
		}},
	}
	data, _ := json.Marshal(exampleConfig)
	file.Write([]byte(data))
}
