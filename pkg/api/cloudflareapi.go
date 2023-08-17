package api

import (
	"context"
	"log"
	"time"

	"github.com/cloudflare/cloudflare-go"
)

// UpdateRecord updates the record with the given recordID
func UpdateRecord(apiToken string, zoneID string, recordID string, recordType RecordType, name string, content string, proxied bool) (*Response, error) {
	clapi, err := cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	zone := cloudflare.ZoneIdentifier(zoneID)
	params := cloudflare.UpdateDNSRecordParams{
		Type:    string(recordType),
		Name:    name,
		Content: content,
		ID:      recordID,
		Proxied: &proxied,
	}
	record, err := clapi.UpdateDNSRecord(ctx, zone, params)
	resp := Response{
		Success:  false,
		Errors:   nil,
		Messages: nil,
		Result: Result{
			ID:         record.ID,
			Type:       record.Type,
			Name:       record.Name,
			Content:    record.Content,
			Proxiable:  record.Proxiable,
			Proxied:    *record.Proxied,
			TTL:        record.TTL,
			Locked:     record.Locked,
			ZoneID:     record.ZoneID,
			ZoneName:   record.ZoneName,
			CreatedOn:  record.CreatedOn,
			ModifiedOn: record.ModifiedOn,
			Data:       record.Data,
			Meta:       record.Meta,
		},
	}
	if err != nil {
		resp.Success = false
	} else {
		resp.Success = true
	}
	return &resp, err
}

// CreateRecord creates a new record
func CreateRecord(apiToken string, zoneID string, recordType RecordType, name string, content string, proxied bool) (*Response, error) {
	clapi, err := cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	ctx := context.TODO()
	zone := cloudflare.ZoneIdentifier(zoneID)
	params := cloudflare.CreateDNSRecordParams{
		CreatedOn:  time.Now(),
		ModifiedOn: time.Now(),
		Type:       string(recordType),
		Name:       name,
		Content:    content,
		Proxied:    &proxied,
	}
	record, err := clapi.CreateDNSRecord(ctx, zone, params)
	resp := Response{
		Success:  false,
		Errors:   nil,
		Messages: nil,
		Result: Result{
			ID:         record.ID,
			Type:       record.Type,
			Name:       record.Name,
			Content:    record.Content,
			Proxiable:  record.Proxiable,
			Proxied:    false,
			TTL:        record.TTL,
			Locked:     record.Locked,
			ZoneID:     record.ZoneID,
			ZoneName:   record.ZoneName,
			CreatedOn:  record.CreatedOn,
			ModifiedOn: record.ModifiedOn,
			Data:       record.Data,
			Meta:       record.Meta,
		},
	}
	if err != nil {
		resp.Success = false
	} else {
		resp.Success = true
	}
	return &resp, err
}

// DeleteRecord deletes record with the given recordID
func DeleteRecord(apiToken string, zoneID string, recordID string) (*Response, error) {
	clapi, err := cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	zone := cloudflare.ZoneIdentifier(zoneID)
	err = clapi.DeleteDNSRecord(ctx, zone, recordID)
	resp := Response{
		Success:  false,
		Errors:   nil,
		Messages: nil,
		Result:   Result{},
	}
	if err != nil {
		resp.Success = false
	} else {
		resp.Success = true
	}
	return &resp, err
}

/*
const url = "https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s"

// UpdateRecord updates the record with the given recordID

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

			body, err = io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			var response Response
			if err = json.Unmarshal(body, &response); err != nil {
				return nil, err
			}
			return &response, nil
		}

// CreateRecord creates a new record

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

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		var response Response
		if err = json.Unmarshal(body, &response); err != nil {
			return nil, err
		}
		return &response, nil
	}

// DeleteRecord deletes record with the given recordID

	func DeleteRecord(apiToken string, zoneID string, recordID string) (*Response, error) {
		resp, err := doAuthorizedRequest(http.MethodDelete, nil, zoneID, recordID, apiToken)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		var response Response
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}
		return &response, nil
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
*/
type requestBody struct {
	RecordType RecordType `json:"type"`
	Name       string     `json:"name"`
	Content    string     `json:"content"`
	TTL        int        `json:"ttl"`
	Proxied    bool       `json:"proxied"`
}

// Response represents the response from the Cloudflare API
type Response struct {
	Success  bool          `json:"success"`  //request was successful
	Errors   []interface{} `json:"errors"`   //potential errors with the request
	Messages []interface{} `json:"messages"` //messages
	Result   Result        `json:"result"`   //See Result
}

// ListedResponse is only applicable to ListRecords as it can
// return more than one result in an Array
type ListedResponse struct {
	Success  bool          `json:"success"`  //request was successful
	Errors   []interface{} `json:"errors"`   //potential errors with the request
	Messages []interface{} `json:"messages"` //messages
	Result   []Result      `json:"result"`   //Array of Result
}

// RecordType represents the type a record can be
type RecordType string

// Possible options for RecordType
const (
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

// Result resembles a record returned from Cloudflare API
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
	Data       interface {
	} `json:"data"`
	Meta interface {
	} `json:"meta"`
}
