package main

import (
	"encoding/json"
	"fmt"
	"github.com/TimoS254/go-dyndns/internal"
	"github.com/TimoS254/go-dyndns/internal/config"
	"github.com/TimoS254/go-dyndns/internal/updater"
	"github.com/TimoS254/go-dyndns/pkg/api"
	"log"
	"os"
	"os/signal"
)

var (
	//GitVersion represents the Git Tag of the current build
	GitVersion string
	//GitBranch represents the Git Branch of the current build
	GitBranch string
	conf      = config.Config{
		IntervalMinutes: 5,
		IntervalSeconds: 0,
		Domains: []config.Domain{{
			DomainName:     "example.com",
			IP4:            true,
			IP6:            true,
			Proxy:          false,
			APIToken:       "yourAPIToken",
			ZoneIdentifier: "yourZoneIdentifier",
		}},
	}
)

func main() {
	log.Printf("Starting go-dyndns")
	log.Printf("Version: %s Branch %s\n", GitVersion, GitBranch)

	//Init Colors
	internal.InitColor()

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

	updater.Update(&conf)
}

func initDomains() {
	for i, domain := range conf.Domains {
		if domain.IP4 {
			ip, _ := api.GetIPv4()
			response, err := api.CreateRecord(domain.APIToken, domain.ZoneIdentifier, api.A, domain.DomainName, ip, domain.Proxy)
			if err != nil {
				log.Printf("Encountered an error while creating A record of Domain %s: %v", domain.DomainName, err)
			}
			if response.Success {
				log.Println("Successfully created A record " + response.Result.Name + " to " + response.Result.Content)
				conf.Domains[i].SetID4(response.Result.ID)
			} else {
				log.Println("Encountered an error while creating " + domain.DomainName + ":")
				fmt.Println(response.Errors)
			}
		}
		if domain.IP6 {
			ip, _ := api.GetIPv6()
			response, err := api.CreateRecord(domain.APIToken, domain.ZoneIdentifier, api.AAAA, domain.DomainName, ip, domain.Proxy)
			if err != nil {
				log.Printf("Encountered an error while creating AAAA record of Domain %s: %v", domain.DomainName, err)
			}
			if response.Success {
				log.Println("Successfully created AAAA record " + response.Result.Name + " to " + response.Result.Content)
				conf.Domains[i].SetID6(response.Result.ID)
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
		config.GenerateConfig(conf)
		log.Println("Created Config from Template!")
		os.Exit(0)
	} else {
		panic(err)
	}
	data, err := os.ReadFile("config.json")
	if err != nil {
		log.Panicf("Could not read Config File: %v", err)
	}
	err = json.Unmarshal(data, &conf)
	if err != nil {
		log.Panicf("Encountered Error unmarshalling json body of config: %v", err)
	}
	log.Println("Loaded Config!")
}

func shutdown() {
	log.Println("Deleting Records and shutting down")
	for _, d := range conf.Domains {
		updater.DeleteRecords(&d)
	}
	return
}
