package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/TimoSLE/go-dyndns/internal/config"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var HttpClient = &http.Client{}

const url = "https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s"

func SetIP(domain config.Domain, recordType string, name string, content string) *Response {
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
		id = domain.GetID4()
	} else if recordType == "AAAA" {
		id = domain.GetID6()
	} else {
		errorResponse := Response{
			Success:  false,
			Errors:   []interface{}{"Wrong Record Type"},
			Messages: []interface{}{"Wrong Record Type"},
		}
		return &errorResponse
	}

	resp := doAuthorizedRequest(http.MethodPut, bytes.NewReader(body), domain.ZoneIdentifier, id, domain.APIToken)

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var response Response
	if err = json.Unmarshal(body, &response); err != nil {
		panic(err)
	}
	return &response
}

func ListRecords(domain config.Domain, forceReqs bool, name string, recordType string) *ListedResponse {
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

	resp := doAuthorizedRequest(http.MethodGet, nil, domain.ZoneIdentifier, s, domain.APIToken)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var response ListedResponse
	if err = json.Unmarshal(body, &response); err != nil {
		panic(err)
	}
	return &response
}

func CreateRecord(domain config.Domain, recordType string, name string, content string) *Response {
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

	resp := doAuthorizedRequest(http.MethodPost, bytes.NewReader(body), domain.ZoneIdentifier, "", domain.APIToken)

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var response Response
	if err = json.Unmarshal(body, &response); err != nil {
		panic(err)
	}
	return &response
}

func DeleteRecords(domain config.Domain) {
	if domain.IP4 {
		resp := doAuthorizedRequest(http.MethodDelete, nil, domain.ZoneIdentifier, domain.GetID4(), domain.APIToken)
		defer resp.Body.Close()
		o, _ := ioutil.ReadAll(resp.Body)
		s := string(o)
		if strings.Contains(s, domain.GetID4()) {
			log.Println("Successfully removed IPv4 Record for " + domain.DomainName)
		}
	}
	if domain.IP6 {
		resp := doAuthorizedRequest(http.MethodDelete, nil, domain.ZoneIdentifier, domain.GetID6(), domain.APIToken)
		defer resp.Body.Close()
		o, _ := ioutil.ReadAll(resp.Body)
		s := string(o)
		if strings.Contains(s, domain.GetID6()) {
			log.Println("Successfully removed IPv6 Record for " + domain.DomainName)
		}
	}
}

func doAuthorizedRequest(method string, body io.Reader, zoneID string, domainID string, apiToken string) *http.Response {
	//Create Request
	req, err := http.NewRequest(method, fmt.Sprintf(url, zoneID, domainID), body)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Bearer "+apiToken)
	req.Header.Set("Content-Type", "application/json")
	req.Close = true

	//Sending Request
	resp, err := HttpClient.Do(req)
	if err != nil {
		panic(err)
	}
	return resp
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
