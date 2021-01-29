package main

import (
	"fmt"
	"log"
	"time"
)

func update() {
	for range time.Tick(time.Minute * time.Duration(config.IntervalMinutes)) {
		for _, domain := range config.Domains {
			if domain.IP4 {
				ip, _ := getIPv4()
				response := setIP(domain, "A", domain.DomainName, ip)
				if response.Success {
					log.Println("Successfully changed A record " + response.Result.Name + " to " + response.Result.Content)
				} else {
					log.Println("Encountered an error while changing " + domain.DomainName + ":")
					fmt.Println(response.Errors)
				}
			}
			if domain.IP6 {
				ip, _ := getIPv6()
				response := setIP(domain, "AAAA", domain.DomainName, ip)
				if response.Success {
					log.Println("Successfully changed A record " + response.Result.Name + " to " + response.Result.Content)
				} else {
					log.Println("Encountered an error while changing " + domain.DomainName + ":")
					fmt.Println(response.Errors)
				}
			}
		}
		httpClient.CloseIdleConnections()
	}
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
