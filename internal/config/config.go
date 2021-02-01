package config

import (
	"encoding/json"
	"log"
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
	id4            string
	id6            string
	APIToken       string
	ZoneIdentifier string
}

func GenerateConfig(config Config) {
	file, err := os.Create("config.json")
	if err != nil {
		log.Panicf("Unexpected Error creating Config File %v", err)
	}
	defer file.Close()
	data, _ := json.MarshalIndent(config, "", "  ")
	_, err = file.Write([]byte(data))
	if err != nil {
		log.Panicf("Unexpected Errro writing Config File %v", err)
	}
}

func (d *Domain) GetID4() string {
	return d.id4
}

func (d *Domain) GetID6() string {
	return d.id6
}

func (d *Domain) SetID4(id string) {
	d.id4 = id
}

func (d *Domain) SetID6(id string) {
	d.id6 = id
}
