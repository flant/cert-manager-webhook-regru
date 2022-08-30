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
	username            string
	password            string
	zone                string
	dumpRequestResponse bool
}

func NewRegruCient(username string, password string, zone string) *RegruClient {
	return &RegruClient{
		username:            username,
		password:            password,
		zone:                zone,
		dumpRequestResponse: false,
	}
}

func (c *RegruClient) getRecords() {
	s := fmt.Sprintf("{\"username\":\"%s\",\"password\":\"%s\",\"domains\":[{\"dname\":\"%s\"}],\"output_content_type\":\"plain\"}", c.username, c.password, c.zone)
	url := fmt.Sprintf("%szone/get_resource_records?input_data=%s&input_format=json", defaultBaseURL, url2.QueryEscape(s))

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

func (c *RegruClient) createTXT() {
	s := fmt.Sprintf("{\"username\":\"%s\",\"password\":\"%s\",\"domains\":[{\"dname\":\"%s\"}],\"subdomain\":\"%s\",\"output_content_type\":\"plain\"}", c.username, c.password, c.zone, domain)
}
