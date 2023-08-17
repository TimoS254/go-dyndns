package api

import (
	"fmt"
	"io"
	"net/http"
)

const ipUrl = "https://api%s.ipify.org"

// HttpClient handles all the http requests of the api
var HttpClient = &http.Client{}

// GetIPv4 returns the current IPv4 of the System in a string
func GetIPv4() (string, error) {
	return getIP("4")
}

// GetIPv6 returns the current IPv6 of the System in a string
func GetIPv6() (string, error) {
	return getIP("6")
}

func getIP(version string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(ipUrl, version), nil)
	if err != nil {
		return "", fmt.Errorf("encountered error creating request %v: %w", req, err)
	}
	req.Close = true
	resp, err := HttpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("encountered error sending request %v: %w", resp, err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("encountered error reading request %v: %w", body, err)
	}
	s := string(body)
	err = resp.Body.Close()
	if err != nil {
		return s, fmt.Errorf("encountered error closing response Body: %w", err)
	}
	resp.Close = true
	return s, nil
}
