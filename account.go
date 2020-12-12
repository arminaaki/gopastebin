package gopastebin

import (
	"bytes"
	"context"
	"encoding/xml"
	"io"
	"log"
	"net/http"
)

const (
	defaultBaseURL  = "https://pastebin.com/"
	accountBasePath = "api/api_login.php"
)

type AccountServiceOp struct {
	client *Client
}

type AccountService interface {
	Create(context.Context, *AccountRequest) (*Account, *Response, error)
	GetUser(context.Context, *AccountUserRequest) ([]AccountUser, *Response, error)
}

// AccountRequest represents a request to create an account
type AccountRequest struct {
	APIDevKey       string `json:"api_dev_key"`
	APIUserName     string `json:"api_user_name"`
	APIUserPassword string `json:"api_user_password"`
}

// Account struct
type Account struct {
	APIUserKey string `json:"api_user_key"`
}

type AccountUserRequest struct {
	APIDevKey  string `json:"api_dev_key,omitempty"`
	APIUserKey string `json:"api_user_key,omitempty"`
	APIOption  string `json:"api_option,omitempty"`
}

type AccountUser struct {
	UserName        string `xml:"user_name,omitempty"`
	UserFormatShort string `xml:"user_format_short,omitempty"`
	UserExpiration  string `xml:"user_expiration,omitempty"`
	UserAvatarURL   string `xml:"user_avatar_url,omitempty"`
	UserPrivate     string `xml:"user_private,omitempty"`
	UserWebsite     string `xml:"user_website,omitempty"`
	UserEmail       string `xml:"user_email,omitempty"`
	UserLocation    string `xml:"user_location,omitempty"`
	UserAccountType string `xml:"user_account_type,omitempty"`
}

// Create a new API KEY given credentials.
func (fw *AccountServiceOp) Create(ctx context.Context, ar *AccountRequest) (*Account, *Response, error) {

	req, err := fw.client.NewRequest(ctx, http.MethodPost, accountBasePath, ar)
	if err != nil {
		return nil, nil, err
	}

	resp, err := fw.client.Do(ctx, req.request)
	if err != nil {
		return nil, resp, err
	}

	return &Account{APIUserKey: resp.response.String()}, resp, nil
}

func (fw *AccountServiceOp) GetUser(ctx context.Context, aur *AccountUserRequest) ([]AccountUser, *Response, error) {
	aur.APIDevKey = fw.client.APIDevKey
	aur.APIUserKey = fw.client.APIUserKey
	aur.APIOption = "userdetails"

	req, err := fw.client.NewRequest(ctx, http.MethodPost, pasteBasePath, aur)
	if err != nil {
		return nil, nil, err
	}

	resp, err := fw.client.Do(ctx, req.request)
	if err != nil {
		return nil, nil, err
	}

	var accountUsers []AccountUser
	d := xml.NewDecoder(bytes.NewBufferString(resp.response.String()))
	for {
		var t AccountUser
		err := d.Decode(&t)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		accountUsers = append(accountUsers, t)
	}
	return accountUsers, resp, nil
}
