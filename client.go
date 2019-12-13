package gopastebin

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"regexp"

	"gopkg.in/resty.v1"
)

// Request wrraper for resty
type Request struct {
	request *resty.Request
}

// Response wrraper for resty
type Response struct {
	response *resty.Response
}

// Client wrraper for resty
type Client struct {
	client     *resty.Client
	BaseURL    *url.URL
	APIUserKey string
	APIDevKey  string

	Account AccountService
	Paste   PasteService
}

// NewClient Creates new client for with base url
func NewClient(accountRequest *AccountRequest) (*Client, error) {
	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{
		client:  resty.New(),
		BaseURL: baseURL,
	}
	c.Account = &AccountServiceOp{client: c}
	c.Paste = &PasteServiceOp{client: c}
	// Auth
	result, _, err := c.Account.Create(nil, accountRequest)
	if err != nil {
		return nil, err
	}
	c.APIUserKey = result.APIUserKey
	c.APIDevKey = accountRequest.APIDevKey

	return c, nil
}

// NewRequest TODO
func (c *Client) NewRequest(ctx context.Context, method string, urlStr string, body interface{}) (*Request, error) {

	request := c.client.NewRequest()
	request.Method = method

	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	request.URL = u.String()

	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		jsonMap := make(map[string]string)
		json.Unmarshal(data, &jsonMap)

		request.SetFormData(jsonMap)
	}

	return &Request{request: request}, nil
}

func (c *Client) Do(ctx context.Context, req *resty.Request) (*Response, error) {
	req.SetContext(ctx)
	resp, err := req.Execute(req.Method, req.URL)
	if err != nil {
		return nil, err
	}

	r, _ := regexp.Compile("invalid")
	if r.MatchString(resp.String()) {
		return &Response{response: resp}, errors.New(resp.String())
	}
	return &Response{response: resp}, nil
}
