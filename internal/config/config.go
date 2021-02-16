package config

import (
	"encoding/json"
	"log"
	"os"
)

//Config consists of all variables saved in the config
type Config struct {
	IntervalMinutes int8
	IntervalSeconds int8
	Domains         []Domain
}

//Domain represents the structure of a domain which should be updated
type Domain struct {
	DomainName     string
	IP4            bool
	IP6            bool
	Proxy          bool
	APIToken       string
	ZoneIdentifier string
	id4            string
	id6            string
}

//GenerateConfig generates the config file from a template
func GenerateConfig(config Config) {
	file, err := os.Create("config.json")
	if err != nil {
		log.Panicf("Unexpected Error creating Config File %v", err)
	}
	defer file.Close()
	data, _ := json.MarshalIndent(config, "", "  ")
	_, err = file.Write([]byte(data))
	if err != nil {
		log.Panicf("Unexpected Error writing Config File %v", err)
	}
}

//GetID4 returns the id of the IPv4 record of a Domain
func (d *Domain) GetID4() string {
	return d.id4
}

//GetID6 returns the id of the IPv6 record of a Domain
func (d *Domain) GetID6() string {
	return d.id6
}

//SetID4 sets the id of the IPv4 record of a Domain
func (d *Domain) SetID4(id string) {
	d.id4 = id
}

//SetID6 sets the id of the IPv6 record of a Domain
func (d *Domain) SetID6(id string) {
	d.id6 = id
}
