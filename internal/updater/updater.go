package updater

import (
	"fmt"
	"github.com/TimoSLE/go-dyndns/internal/config"
	"github.com/TimoSLE/go-dyndns/pkg/api"
	"log"
	"time"
)

func Update(conf config.Config) {

	for range time.Tick(time.Minute * time.Duration(conf.IntervalMinutes)) {
		for _, domain := range conf.Domains {
			if domain.IP4 {
				ip, _ := api.GetIPv4()
				response := api.SetIP(domain, "A", domain.DomainName, ip)
				if response.Success {
					log.Println("Successfully changed A record " + response.Result.Name + " to " + response.Result.Content)
				} else {
					log.Println("Encountered an error while changing " + domain.DomainName + ":")
					fmt.Println(response.Errors)
				}
			}
			if domain.IP6 {
				ip, _ := api.GetIPv6()
				response := api.SetIP(domain, "AAAA", domain.DomainName, ip)
				if response.Success {
					log.Println("Successfully changed AAAA record " + response.Result.Name + " to " + response.Result.Content)
				} else {
					log.Println("Encountered an error while changing " + domain.DomainName + ":")
					fmt.Println(response.Errors)
				}
			}
		}
		api.HttpClient.CloseIdleConnections()
	}
}
