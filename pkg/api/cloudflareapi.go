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

func SetIP(apiToken string, zoneID string, domainID string, recordType string, name string, content string) *Response {
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

	resp := doAuthorizedRequest(http.MethodPut, bytes.NewReader(body), zoneID, domainID, apiToken)

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

func ListRecords(apiToken string, zoneID string, forceReqs bool, name string, recordType string) *ListedResponse {
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

	resp := doAuthorizedRequest(http.MethodGet, nil, zoneID, s, apiToken)

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

func CreateRecord(apiToken string, zoneID string, recordType string, name string, content string) *Response {
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

	resp := doAuthorizedRequest(http.MethodPost, bytes.NewReader(body), zoneID, "", apiToken)

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

func DeleteRecord(apiToken string, zoneID string, domainID string) *Result {
	resp := doAuthorizedRequest(http.MethodDelete, nil, zoneID, domainID, apiToken)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var result Result
	json.Unmarshal(body, &result)
	return &result
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
