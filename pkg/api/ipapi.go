package api

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

func GetIPv4() (string, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api4.publicip.xyz", nil)
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

func GetIPv6() (string, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api6.publicip.xyz", nil)
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
