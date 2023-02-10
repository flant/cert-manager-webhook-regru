package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	url2 "net/url"
)

const (
	defaultBaseURL = "https://api.reg.ru/api/regru2/"
)

type RegruClient struct {
	username string
	password string
	zone     string
}

func NewRegruClient(username string, password string, zone string) *RegruClient {
	return &RegruClient{
		username: username,
		password: password,
		zone:     zone,
	}
}

func (c *RegruClient) getRecords() error {
	s := fmt.Sprintf("{\"username\":\"%s\",\"password\":\"%s\",\"domains\":[{\"dname\":\"%s\"}],\"output_content_type\":\"plain\"}", c.username, c.password, c.zone)
	url := fmt.Sprintf("%szone/get_resource_records?input_data=%s&input_format=json", defaultBaseURL, url2.QueryEscape(s))

	fmt.Println("Get TXT Query:", url)
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to make GET request: %v", err)
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		return errors.New(fmt.Sprintf("response failed with status code: %d and body: %s", res.StatusCode, body))
	}
	if err != nil {
		return fmt.Errorf("failed to ready response")
	}

	fmt.Printf("Get TXT success. Response body: %s", body)

	return nil
}

func (c *RegruClient) createTXT(domain string, value string) error {
	s := fmt.Sprintf("{\"username\":\"%s\",\"password\":\"%s\",\"domains\":[{\"dname\":\"%s\"}],\"subdomain\":\"%s\",\"text\":\"%s\",\"output_content_type\":\"plain\"}", c.username, c.password, c.zone, domain, value)
	url := fmt.Sprintf("%szone/add_txt?input_data=%s&input_format=json", defaultBaseURL, url2.QueryEscape(s))

	fmt.Println("Create TXT Query:", url)
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to make GET request: %v", err)
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		return errors.New(fmt.Sprintf("response failed with status code: %d and body: %s", res.StatusCode, body))
	}
	if err != nil {
		return fmt.Errorf("failed to ready response")
	}

	fmt.Printf("Create TXT success. Response body: %s", body)

	return nil
}

func (c *RegruClient) deleteTXT(domain string, value string) error {
	s := fmt.Sprintf("{\"username\":\"%s\",\"password\":\"%s\",\"domains\":[{\"dname\":\"%s\"}],\"subdomain\":\"%s\",\"content\":\"%s\",\"record_type\":\"TXT\",\"output_content_type\":\"plain\"}", c.username, c.password, c.zone, domain, value)
	url := fmt.Sprintf("%szone/remove_record?input_data=%s&input_format=json", defaultBaseURL, url2.QueryEscape(s))

	fmt.Println("Delete TXT Query:", url)
	res, err := http.Get(url)

	if err != nil {
		return fmt.Errorf("failed to make GET request: %v", err)
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		return errors.New(fmt.Sprintf("response failed with status code: %d and body: %s", res.StatusCode, body))
	}
	if err != nil {
		return fmt.Errorf("failed to ready response")
	}

	fmt.Printf("Delete TXT success. Response body: %s", body)
	return nil
}
