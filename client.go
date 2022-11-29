package main

import (
	"fmt"
	"io/ioutil"
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

func NewRegruCient(username string, password string, zone string) *RegruClient {
	return &RegruClient{
		username: username,
		password: password,
		zone:     zone,
	}
}

func (c *RegruClient) getRecords() {
	s := fmt.Sprintf("{\"username\":\"%s\",\"password\":\"%s\",\"domains\":[{\"dname\":\"%s\"}],\"output_content_type\":\"plain\"}", c.username, c.password, c.zone)
	url := fmt.Sprintf("%szone/get_resource_records?input_data=%s&input_format=json", defaultBaseURL, url2.QueryEscape(s))

	fmt.Println("Query:", url)
	req, err := http.Get(url)

	if err != nil {
		fmt.Println(err)
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))

}

func (c *RegruClient) createTXT(domain string, value string) error {
	s := fmt.Sprintf("{\"username\":\"%s\",\"password\":\"%s\",\"domains\":[{\"dname\":\"%s\"}],\"subdomain\":\"%s\",\"text\":\"%s\",\"output_content_type\":\"plain\"}", c.username, c.password, c.zone, domain, value)
	url := fmt.Sprintf("%szone/add_txt?input_data=%s&input_format=json", defaultBaseURL, url2.QueryEscape(s))

	fmt.Println("Query:", url)
	req, err := http.Get(url)

	if err != nil {
		return fmt.Errorf("failed creating TXT record: %v", err)
	}

	body, _ := ioutil.ReadAll(req.Body)
	fmt.Println(string(body))

	if err != nil {
		fmt.Sprintf("Created TXT record: %s", err)
	} else {
		fmt.Sprintf("Created TXT record: %s", body)
	}
	return nil
}

func (c *RegruClient) deleteTXT(domain string, value string) error {
	s := fmt.Sprintf("{\"username\":\"%s\",\"password\":\"%s\",\"domains\":[{\"dname\":\"%s\"}],\"subdomain\":\"%s\",\"content\":\"%s\",\"record_type\":\"TXT\",\"output_content_type\":\"plain\"}", c.username, c.password, c.zone, domain, value)
	url := fmt.Sprintf("%szone/remove_record?input_data=%s&input_format=json", defaultBaseURL, url2.QueryEscape(s))

	fmt.Println("Query:", url)
	req, err := http.Get(url)

	if err != nil {
		return fmt.Errorf("failed creating TXT record: %v", err)
	}

	body, _ := ioutil.ReadAll(req.Body)
	fmt.Sprintf("Created TXT record: %s", body)
	return nil
}
