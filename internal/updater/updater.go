package updater

import (
	"github.com/TimoSLE/go-dyndns/internal/config"
	"github.com/TimoSLE/go-dyndns/pkg/api"
	"log"
	"strings"
	"time"
)

//Update Starts the ticker for updating domains
func Update(conf *config.Config) {
	var delay = time.Duration(0)
	if conf.IntervalMinutes != 0 {
		delay = delay + time.Minute*time.Duration(conf.IntervalMinutes)
	}
	if conf.IntervalSeconds != 0 {
		delay = delay + time.Second*time.Duration(conf.IntervalSeconds)
	}
	if delay.Seconds() < 1 {
		delay = time.Second * 1
	}
	for range time.Tick(delay) {
		for _, domain := range conf.Domains {
			go UpdateDomain(&domain)
		}
		api.HttpClient.CloseIdleConnections()
	}
}

//UpdateDomain handles the update of a specific config.Domain
func UpdateDomain(domain *config.Domain) {
	if domain.IP4 {
		ip, _ := api.GetIPv4()
		response, err := api.UpdateRecord(domain.APIToken, domain.ZoneIdentifier, domain.GetID4(), api.A, domain.DomainName, ip, domain.Proxy)
		if err != nil {
			log.Printf("Encountered an error while updating A record of Domain %s: %v", domain.DomainName, err)
		}
		if response.Success {
			log.Println("Successfully changed A record " + response.Result.Name + " to " + response.Result.Content)
		} else {
			log.Printf("Encountered an error while changing %s: %v ", domain.DomainName, response.Errors)
		}
	}
	if domain.IP6 {
		ip, _ := api.GetIPv6()
		response, err := api.UpdateRecord(domain.APIToken, domain.ZoneIdentifier, domain.GetID6(), api.AAAA, domain.DomainName, ip, domain.Proxy)
		if err != nil {
			log.Printf("Encountered an error while updating AAAA record of Domain %s: %v", domain.DomainName, err)
		}
		if response.Success {
			log.Println("Successfully changed AAAA record " + response.Result.Name + " to " + response.Result.Content)
		} else {
			log.Printf("Encountered an error while changing %s: %v ", domain.DomainName, response.Errors)
		}
	}
}

//DeleteRecords deletes all records of a specific Domain
func DeleteRecords(domain *config.Domain) {
	if domain.IP4 {
		response, err := api.DeleteRecord(domain.APIToken, domain.ZoneIdentifier, domain.GetID4())
		if err != nil {
			log.Printf("Encountered an error while deleting A record of Domain %s: %v", domain.DomainName, err)
		}
		if strings.Contains(response.ID, domain.GetID4()) {
			log.Println("Successfully removed IPv4 Record for " + domain.DomainName)
		}
	}
	if domain.IP6 {
		response, err := api.DeleteRecord(domain.APIToken, domain.ZoneIdentifier, domain.GetID6())
		if err != nil {
			log.Printf("Encountered an error while deleting AAAA record of Domain %s: %v", domain.DomainName, err)
		}
		if strings.Contains(response.ID, domain.GetID6()) {
			log.Println("Successfully removed IPv6 Record for " + domain.DomainName)
		}
	}
}
