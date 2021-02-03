package updater

import (
	"fmt"
	"github.com/TimoSLE/go-dyndns/internal/config"
	"github.com/TimoSLE/go-dyndns/pkg/api"
	"log"
	"strings"
	"time"
)

func Update(conf *config.Config) {

	for range time.Tick(time.Minute * time.Duration(conf.IntervalMinutes)) {
		for _, domain := range conf.Domains {
			go UpdateDomain(&domain)
		}
		api.HttpClient.CloseIdleConnections()
	}
}

func UpdateDomain(domain *config.Domain) {
	if domain.IP4 {
		ip, _ := api.GetIPv4()
		response := api.SetIP(domain.APIToken, domain.ZoneIdentifier, domain.GetID4(), "A", domain.DomainName, ip)
		if response.Success {
			log.Println("Successfully changed A record " + response.Result.Name + " to " + response.Result.Content)
		} else {
			log.Println("Encountered an error while changing " + domain.DomainName + ":")
			fmt.Println(response.Errors)
		}
	}
	if domain.IP6 {
		ip, _ := api.GetIPv6()
		response := api.SetIP(domain.APIToken, domain.ZoneIdentifier, domain.GetID6(), "AAAA", domain.DomainName, ip)
		if response.Success {
			log.Println("Successfully changed AAAA record " + response.Result.Name + " to " + response.Result.Content)
		} else {
			log.Println("Encountered an error while changing " + domain.DomainName + ":")
			fmt.Println(response.Errors)
		}
	}
}

func DeleteRecords(domain *config.Domain) {
	if domain.IP4 {
		res := api.DeleteRecord(domain.APIToken, domain.ZoneIdentifier, domain.GetID4())
		if strings.Contains(res.ID, domain.GetID4()) {
			log.Println("Successfully removed IPv4 Record for " + domain.DomainName)
		}
	}
	if domain.IP6 {
		res := api.DeleteRecord(domain.APIToken, domain.ZoneIdentifier, domain.GetID6())
		if strings.Contains(res.ID, domain.GetID6()) {
			log.Println("Successfully removed IPv6 Record for " + domain.DomainName)
		}
	}
}
