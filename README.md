# go-dyndns [![Go Report Card](https://goreportcard.com/badge/github.com/TimoSLE/go-dyndns)](https://goreportcard.com/report/github.com/TimoSLE/go-dyndns) ![Build](https://github.com/TimoSLE/go-dyndns/workflows/Auto%20Build%20and%20Release/badge.svg)



## Building
You can download an executable in the releases Section or build it with make using
> make build

If you are not using make you can just build it with
>go build ./cmd/go-dyndns/main.go

## Usage
On First Execution the Program will generate a config.json file which you have to modify according to your Usecase.  
  
The Fol

Config Key | Datatype | Example | Explanation
------------ | ------------- | ------------- | -------------
IntervalMinutes | int8 | 5 | Update Interval in Minutes
Domains | Array | / | Array of Domains which should be updated

Each of the Domains in the Domain Array should consist of the following variables

Key | Datatype | Example | Explanation
------------ | ------------- | ------------- | -------------
DomainName | string | "example.com" | Name of the Record
IP4 | bool | true | Whether a A record should be created/updated
IP6 | bool | true | Whether a AAAA record should be created/updated
Proxy | bool | false | Should the domain be proxied through Cloudflare CDN
APIToken | string | "" | Your Cloudflare API token (needs DNS Edit Permissions)
ZoneIdentifier | string | "" | Zone ID of the Zone the record should be located



## Contributing
Every Contribution is welcome, feel free to open pull requests and/or suggest changes in issues