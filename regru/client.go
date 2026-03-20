package regru

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultBaseURL = "https://api.reg.ru/api/regru2/"

// defaultClient is the single reg.ru API client used by the solver.
var regruClient *Client

// Client is the reg.ru API client used to manage DNS records.
type Client struct {
	username   string
	password   string
	baseURL    *url.URL
	HTTPClient *http.Client
}

// NewClient creates a new reg.ru API client with the provided credentials.
func NewClient(username, password string) *Client {
	baseURL, err := url.Parse(defaultBaseURL)
	if err != nil {
		panic("regru: invalid default base URL: " + err.Error())
	}
	return &Client{
		username:   username,
		password:   password,
		baseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// InitClient sets the package-level default Client. Must be called once before NewSolver.
func InitClient(username, password string) {
	regruClient = NewClient(username, password)
}

// addTXTRecord adds a TXT DNS record for the given domain and subdomain via the reg.ru API.
func (c *Client) addTXTRecord(ctx context.Context, domain, subdomain, txtContent string) error {
	request := AddTxtRequest{
		Domains:           []Domain{{DName: domain}},
		SubDomain:         subdomain,
		Text:              txtContent,
		OutputContentType: "json",
	}

	err := c.doRequest(ctx, request, "zone/add_txt")
	if err != nil {
		return fmt.Errorf("failed to add TXT record: %v", err)
	}

	return nil
}

// deleteTXTRecord removes a TXT DNS record for the given domain and subdomain via the reg.ru API.
func (c *Client) deleteTXTRecord(ctx context.Context, domain, subdomain, txtContent string) error {
	request := RemoveRecordRequest{
		Domains:           []Domain{{DName: domain}},
		SubDomain:         subdomain,
		Content:           txtContent,
		RecordType:        "TXT",
		OutputContentType: "json",
	}

	err := c.doRequest(ctx, request, "zone/remove_record")
	if err != nil {
		return fmt.Errorf("failed to delete TXT record: %v", err)
	}
	return nil
}

// doRequest serializes the request as JSON, sends it to the given reg.ru API endpoint,
// and logs the response. Credentials are passed as form fields per the reg.ru API spec.
func (c *Client) doRequest(ctx context.Context, request any, apiEndpoint string) error {
	endpoint := c.baseURL.JoinPath(apiEndpoint)

	inputData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	data := url.Values{}
	data.Set("username", c.username)
	data.Set("password", c.password)
	data.Set("input_format", "json")
	data.Set("output_format", "json")
	data.Set("io_encoding", "utf8")
	data.Set("input_data", string(inputData))
	data.Set("show_input_params", "0")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	slog.Default().Info("sending API request", "request", request, "path", apiEndpoint)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, errRead := io.ReadAll(resp.Body)
		if errRead != nil {
			return fmt.Errorf("API request failed: HTTP %s (body read error: %w)", resp.Status, errRead)
		}
		return fmt.Errorf("API request failed: HTTP %s: %s", resp.Status, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResp apiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}
	if apiResp.Result != "success" {
		msg := apiResp.ErrorText
		if msg == "" {
			msg = apiResp.Result
		}
		return fmt.Errorf("API request failed: %s", msg)
	}

	var jsonResponse interface{}
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	prettyJSON, err := json.MarshalIndent(jsonResponse, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	slog.Default().Info("received API response",
		"path", apiEndpoint,
		"status_code", resp.StatusCode,
		"response_body", prettyJSON,
	)

	return nil
}
