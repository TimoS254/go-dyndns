package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

var HttpClient = &http.Client{}

const url = "https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s"

//UpdateRecord updates the record with the given recordID
func UpdateRecord(apiToken string, zoneID string, recordID string, recordType RecordType, name string, content string, proxied bool) (*Response, error) {
	//Creating request Body
	request := requestBody{
		RecordType: recordType,
		Name:       name,
		Content:    content,
		TTL:        1,
		Proxied:    proxied,
	}
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	resp, err := doAuthorizedRequest(http.MethodPut, bytes.NewReader(body), zoneID, recordID, apiToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response Response
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

//ListRecords list all records which fit the given parameters
func ListRecords(apiToken string, zoneID string, forceReqs bool, name string, recordType RecordType) (*ListedResponse, error) {
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
		s = s + "&type=" + string(recordType)
	}
	s = s + "&per_page=100"

	resp, err := doAuthorizedRequest(http.MethodGet, nil, zoneID, s, apiToken)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var response ListedResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

//CreateRecord creates a new record
func CreateRecord(apiToken string, zoneID string, recordType RecordType, name string, content string, proxied bool) (*Response, error) {
	//Creating Json request Body
	request := requestBody{
		RecordType: recordType,
		Name:       name,
		Content:    content,
		TTL:        1,
		Proxied:    proxied,
	}
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	resp, err := doAuthorizedRequest(http.MethodPost, bytes.NewReader(body), zoneID, "", apiToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response Response
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

//DeleteRecord deletes record with the given recordID
func DeleteRecord(apiToken string, zoneID string, recordID string) (*Result, error) {
	resp, err := doAuthorizedRequest(http.MethodDelete, nil, zoneID, recordID, apiToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var result Result
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func doAuthorizedRequest(method string, body io.Reader, zoneID string, domainID string, apiToken string) (*http.Response, error) {
	//Create Request
	req, err := http.NewRequest(method, fmt.Sprintf(url, zoneID, domainID), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+apiToken)
	req.Header.Set("Content-Type", "application/json")
	req.Close = true

	//Sending Request
	resp, err := HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type requestBody struct {
	RecordType RecordType `json:"type"`
	Name       string     `json:"name"`
	Content    string     `json:"content"`
	TTL        int        `json:"ttl"`
	Proxied    bool       `json:"proxied"`
}

//Response represents the response from the Cloudflare API
type Response struct {
	Success  bool          `json:"success"`  //request was successful
	Errors   []interface{} `json:"errors"`   //potential errors with the request
	Messages []interface{} `json:"messages"` //messages
	Result   Result        `json:"result"`   //See Result
}

//ListedResponse is only applicable to ListRecords as it can
//return more than one result in an Array
type ListedResponse struct {
	Success  bool          `json:"success"`  //request was successful
	Errors   []interface{} `json:"errors"`   //potential errors with the request
	Messages []interface{} `json:"messages"` //messages
	Result   []Result      `json:"result"`   //Array of Result
}

//RecordType
type RecordType string

const (
	//Possible options for RecordType
	A     RecordType = "A"
	AAAA  RecordType = "AAAA"
	TXT   RecordType = "TXT"
	CNAME RecordType = "CNAME"
	HTTPS RecordType = "HTTPS"
	SRV   RecordType = "SRV"
	LOC   RecordType = "LOC"
	MX    RecordType = "MX"
	NS    RecordType = "NS"
	SPF   RecordType = "SPF"
)

//Result resembles a record returned from Cloudflare API
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
