package api

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const IPURL = "https://api%s.publicip.xyz"

func GetIPv4() (string, error) {
	return getIP("4")
}

func GetIPv6() (string, error) {
	return getIP("6")
}

func getIP(version string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(IPURL, version), nil)
	if err != nil {
		return "", errors.New("encountered an Error while creating request")
	}
	req.Close = true
	resp, err := HttpClient.Do(req)
	if err != nil {
		return "", errors.New("encountered an Error while sending request")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("encountered an Error while reading response")
	}
	s := string(body)
	err = resp.Body.Close()
	if err != nil {
		log.Printf("Couldnt Close Request Body %v", err)
	}
	resp.Close = true
	return s, nil
}
