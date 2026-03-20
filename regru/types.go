package regru

// AddTxtRequest is the representation of the payload of a request to add a TXT record.
type AddTxtRequest struct {
	Domains           []Domain `json:"domains,omitempty"`
	SubDomain         string   `json:"subdomain,omitempty"`
	Text              string   `json:"text,omitempty"`
	OutputContentType string   `json:"output_content_type,omitempty"`
}

// RemoveRecordRequest is the representation of the payload of a request to remove a record.
type RemoveRecordRequest struct {
	Domains           []Domain `json:"domains,omitempty"`
	SubDomain         string   `json:"subdomain,omitempty"`
	Content           string   `json:"content,omitempty"`
	RecordType        string   `json:"record_type,omitempty"`
	OutputContentType string   `json:"output_content_type,omitempty"`
}

// Domain represents a domain entry used in reg.ru API requests.
type Domain struct {
	DName string `json:"dname"`
}

// apiResponse is the minimal structure for reg.ru API JSON response to detect errors.
type apiResponse struct {
	Result    string `json:"result"`
	ErrorText string `json:"error_text"`
}
