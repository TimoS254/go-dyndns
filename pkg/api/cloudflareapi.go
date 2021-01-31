package api

import (
	"bytes"
	"encoding/json"
	"go-dyndns/internal/config"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var HttpClient = &http.Client{}

func SetIP(domain config.Domain, recordType string, name string, content string) Response {
	//Creating Request Body
	request := SetRequestBody{
		RecordType: recordType,
		Name:       name,
		Content:    content,
		TTL:        1,
		Proxied:    false,
	}
	body, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}

	id := ""
	if recordType == "A" {
		id = config.GetID4(domain)
	} else if recordType == "AAAA" {
		id = config.GetID6(domain)
	} else {
		//TODO Error Handling, wrong record type
	}

	//Creating Request and Setting Headers
	req, err := http.NewRequest(http.MethodPut, "https://api.cloudflare.com/client/v4/zones/"+domain.ZoneIdentifier+"/dns_records/"+id, bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Bearer "+domain.APIToken)
	req.Header.Set("Content-Type", "application/json")
	req.Close = true

	//Sending Request
	resp, err := HttpClient.Do(req)
	if err != nil {
		panic(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var response Response
	if err = json.Unmarshal(body, &response); err != nil {
		panic(err)
	}
	return response
}

func listRecords(domain config.Domain, forceReqs bool, name string, recordType string) ListedResponse {
	s := "?"
	temp := "any"
	if forceReqs {
		temp = "all"
	}
	s = s + "match=" + temp
	if name != "" {
		s = s + "&name=" + name
	}
	if recordType != "" {
		s = s + "&type=" + recordType
	}

	//Create Request, Set Headers
	req, err := http.NewRequest(http.MethodGet, "https://api.cloudflare.com/client/v4/zones/"+domain.ZoneIdentifier+"/dns_records"+s, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Bearer "+domain.APIToken)
	req.Header.Set("Content-Type", "application/json")
	req.Close = true

	resp, err := HttpClient.Do(req)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var response ListedResponse
	if err = json.Unmarshal(body, &response); err != nil {
		panic(err)
	}
	return response
}

func CreateRecord(domain config.Domain, recordType string, name string, content string) Response {
	//Creating Json Request Body
	request := SetRequestBody{
		RecordType: recordType,
		Name:       name,
		Content:    content,
		TTL:        1,
		Proxied:    false,
	}
	body, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.cloudflare.com/client/v4/zones/"+domain.ZoneIdentifier+"/dns_records", bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Bearer "+domain.APIToken)
	req.Header.Set("Content-Type", "application/json")
	req.Close = true

	//Sending Request
	resp, err := HttpClient.Do(req)
	if err != nil {
		panic(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var response Response
	if err = json.Unmarshal(body, &response); err != nil {
		panic(err)
	}
	return response
}

func DeleteRecords(domain config.Domain) {
	if domain.IP4 {
		req, err := http.NewRequest(http.MethodDelete, "https://api.cloudflare.com/client/v4/zones/"+domain.ZoneIdentifier+"/dns_records/"+config.GetID4(domain), nil)
		if err != nil {
			panic(err)
		}
		req.Header.Set("Authorization", "Bearer "+domain.APIToken)
		req.Header.Set("Content-Type", "application/json")
		req.Close = true
		resp, err := HttpClient.Do(req)
		defer resp.Body.Close()
		o, _ := ioutil.ReadAll(resp.Body)
		s := string(o)
		if strings.Contains(s, config.GetID4(domain)) {
			log.Println("Successfully removed IPv4 Record for " + domain.DomainName)
		}
	}
	if domain.IP6 {
		req, err := http.NewRequest(http.MethodDelete, "https://api.cloudflare.com/client/v4/zones/"+domain.ZoneIdentifier+"/dns_records/"+config.GetID6(domain), nil)
		if err != nil {
			panic(err)
		}
		req.Header.Set("Authorization", "Bearer "+domain.APIToken)
		req.Header.Set("Content-Type", "application/json")
		req.Close = true
		resp, err := HttpClient.Do(req)
		defer resp.Body.Close()
		o, _ := ioutil.ReadAll(resp.Body)
		s := string(o)
		if strings.Contains(s, config.GetID6(domain)) {
			log.Println("Successfully removed IPv6 Record for " + domain.DomainName)
		}
	}
}

type SetRequestBody struct {
	RecordType string `json:"type"`
	Name       string `json:"name"`
	Content    string `json:"content"`
	TTL        int    `json:"ttl"`
	Proxied    bool   `json:"proxied"`
}

type Response struct {
	Success  bool          `json:"success"`
	Errors   []interface{} `json:"errors"`
	Messages []interface{} `json:"messages"`
	Result   Result        `json:"result"`
}

type ListedResponse struct {
	Success  bool          `json:"success"`
	Errors   []interface{} `json:"errors"`
	Messages []interface{} `json:"messages"`
	Result   []Result      `json:"result"`
}

type Result struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	Name       string    `json:"name"`
	Content    string    `json:"content"`
	Proxiable  bool      `json:"proxiable"`
	Proxied    bool      `json:"proxied"`
	TTL        int       `json:"ttl"`
	Locked     bool      `json:"locked"`
	ZoneID     string    `json:"zone_id"`
	ZoneName   string    `json:"zone_name"`
	CreatedOn  time.Time `json:"created_on"`
	ModifiedOn time.Time `json:"modified_on"`
	Data       struct {
	} `json:"data"`
	Meta struct {
		AutoAdded bool   `json:"auto_added"`
		Source    string `json:"source"`
	} `json:"meta"`
}
