package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
)

var config Config = Config{
	IntervalMinutes: 5,
	Domains: []Domain{{
		DomainName:     "example.com",
		IP4:            true,
		IP6:            true,
		APIToken:       "yourAPIToken",
		ZoneIdentifier: "yourZoneIdentifier",
	}},
}

func main() {
	log.Println("Starting go-dyndns")
	//Initialize Config
	initConfig()

	//Init Domains
	initDomains()

	//OS Interrupt
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	signal.Notify(interrupt, os.Kill)
	go func() {
		<-interrupt
		defer os.Exit(1)
		shutdown()
		return
	}()

	update()
}

func initDomains() {
	for i, domain := range config.Domains {
		if domain.IP4 {
			ip, _ := getIPv4()
			response := createRecord(domain, "A", domain.DomainName, ip)
			if response.Success {
				log.Println("Successfully created A record " + response.Result.Name + " to " + response.Result.Content)
				config.Domains[i].lastID4 = response.Result.ID
			} else {
				log.Println("Encountered an error while creating " + domain.DomainName + ":")
				fmt.Println(response.Errors)
			}
		}
		if domain.IP6 {
			ip, _ := getIPv6()
			response := createRecord(domain, "AAAA", domain.DomainName, ip)
			if response.Success {
				log.Println("Successfully created AAAA record " + response.Result.Name + " to " + response.Result.Content)
				config.Domains[i].lastID6 = response.Result.ID
			} else {
				log.Println("Encountered an error while creating " + domain.DomainName + ":")
				fmt.Println(response.Errors)
			}
		}
	}
}

func initConfig() {
	log.Println("Initializing Config...")
	//Check if Config already exists
	if _, err := os.Stat("config.json"); err == nil {
		//Continue with Initialization
		log.Println("Loading Config...")
	} else if os.IsNotExist(err) {
		log.Println("Creating Config from Template...")
		generateConfig()
		log.Println("Created Config from Template!")
		os.Exit(0)
	} else {
		panic(err)
	}
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &config)
	log.Println("Loaded Config!")
}

func shutdown() {
	log.Println("Deleting Records and shutting down")
	for _, d := range config.Domains {
		deleteRecords(d)
	}
	return
}
