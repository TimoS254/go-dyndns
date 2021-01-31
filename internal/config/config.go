package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	IntervalMinutes int8
	Domains         []Domain
}

type Domain struct {
	DomainName     string
	IP4            bool
	IP6            bool
	LastID4        string
	LastID6        string
	APIToken       string
	ZoneIdentifier string
}

func GenerateConfig(config Config) {
	file, _ := os.Create("config.json")
	defer file.Close()
	data, _ := json.MarshalIndent(config, "", "  ")
	file.Write([]byte(data))
}
