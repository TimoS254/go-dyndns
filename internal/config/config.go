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
	lastID4        string
	lastID6        string
	APIToken       string
	ZoneIdentifier string
}

func GetID4(domain Domain) string {
	return domain.lastID4
}

func SetID4(domain Domain, id string) {
	domain.lastID4 = id
}

func GetID6(domain Domain) string {
	return domain.lastID6
}

func SetID6(domain Domain, id string) {
	domain.lastID6 = id
}

func GenerateConfig(config Config) {
	file, _ := os.Create("config.json")
	defer file.Close()
	data, _ := json.MarshalIndent(config, "", "  ")
	file.Write([]byte(data))
}
