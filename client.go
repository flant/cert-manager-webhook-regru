package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
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
	apiURL := fmt.Sprintf("%szone/get_resource_records", defaultBaseURL)
	inputData := fmt.Sprintf("{\"domains\":[{\"dname\":\"%s\"}],\"password\":\"%s\",\"username\":\"%s\"}", c.zone, c.password, c.username)
	return sendPOST(apiURL, inputData, *c)
}

func (c *RegruClient) createTXT(domain string, value string) error {
	apiURL := fmt.Sprintf("%szone/add_txt", defaultBaseURL)
	inputData := fmt.Sprintf("{\"domains\":[{\"dname\":\"%s\"}],\"password\":\"%s\",\"subdomain\":\"%s\",\"text\":\"%s\",\"username\":\"%s\"}", c.zone, c.password, domain, value, c.username)
	return sendPOST(apiURL, inputData, *c)
}

func (c *RegruClient) deleteTXT(domain string, value string) error {
	apiURL := fmt.Sprintf("%szone/remove_record", defaultBaseURL)
	inputData := fmt.Sprintf("{\"content\":\"%s\",\"domains\":[{\"dname\":\"%s\"}],\"password\":\"%s\",\"record_type\":\"TXT\",\"subdomain\":\"%s\",\"username\":\"%s\"}", value, c.zone, c.password, domain, c.username)
	return sendPOST(apiURL, inputData, *c)
}

func sendPOST(apiURL string, inputData string, c RegruClient) error {
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	writer.WriteField("input_format", "json")
	writer.WriteField("output_format", "json")
	writer.WriteField("io_encoding", "utf8")
	writer.WriteField("input_data", inputData)
	writer.WriteField("show_input_params", "0")
	writer.WriteField("username", c.username)
	writer.WriteField("password", c.password)
	writer.Close()

	// Perform the POST request
	req, err := http.NewRequest("POST", apiURL, &b)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Perform the POST request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make POST request: %v", err)
	}
	defer res.Body.Close()

	// Check for non-success status code
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body) // Ignore error for brevity
		return fmt.Errorf("response failed with status code: %d and body: %s", res.StatusCode, body)
	}

	// Read and output the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	// Print the response body as formatted JSON
	var jsonResponse interface{}
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		return fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	// Marshal the jsonResponse with indentation for pretty printing
	prettyJSON, err := json.MarshalIndent(jsonResponse, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	fmt.Printf("Response body: %s\n", prettyJSON)
	return nil
}
