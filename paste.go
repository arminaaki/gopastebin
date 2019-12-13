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
	pasteBasePath    = "api/api_post.php"
	rawPasteBasePath = "api/api_raw.php"

	APIOptionPublic   = "0"
	APIOptionUnlisted = "1"
	APIOptionPrivate  = "2"
)

type PasteServiceOp struct {
	client *Client
}

type PasteService interface {
	Create(context.Context, *PasteRequest) (*Paste, *Response, error)
	Delete(context.Context, *PasteDeleteRequest) (*Response, error)
	List(context.Context, *PasteListRequest) ([]Paste, *Response, error)
	GetRaw(ctx context.Context, rpr *PasteGetRawRequest) ([]byte, *Response, error)
}

type PasteRequest struct {
	APIDevKey  string `json:"api_dev_key,omitempty"`
	APIUserKey string `json:"api_user_key,omitempty"`
	APIOption  string `json:"api_option,omitempty"`

	APIPasteCode string `json:"api_paste_code,omitempty"`
	// Optional Params:
	APIPasteName       string `json:"api_paste_name,omitempty"`
	APIPasteFormat     string `json:"api_paste_format,omitempty"`
	APIPastePrivate    string `json:"api_paste_private,omitempty"`
	APIPasteExpireDate string `json:"api_paste_expire_date,omitempty"`
}

type PasteDeleteRequest struct {
	APIDevKey  string `json:"api_dev_key,omitempty"`
	APIUserKey string `json:"api_user_key,omitempty"`
	APIOption  string `json:"api_option,omitempty"`

	APIPasteKey string `json:"api_paste_key,omitempty"`
}

type PasteListRequest struct {
	APIDevKey  string `json:"api_dev_key,omitempty"`
	APIUserKey string `json:"api_user_key,omitempty"`
	APIOption  string `json:"api_option,omitempty"`

	APIResultsLimit string `json:"api_results_limit,omitempty"`
}

type PasteGetRawRequest struct {
	APIDevKey  string `json:"api_dev_key,omitempty"`
	APIUserKey string `json:"api_user_key,omitempty"`
	APIOption  string `json:"api_option,omitempty"`

	PasteKey string `json:"api_paste_key,omitempty"`
}
type Paste struct {
	PasteURL         string `xml:"paste_url"`
	PasteKey         string `xml:"paste_key"`
	PasteDate        string `xml:"paste_date"`
	PasteTitle       string `xml:"paste_title"`
	PasteSize        string `xml:"paste_size"`
	PasteExpireDate  string `xml:"paste_expire_date"`
	PastePrivate     string `xml:"paste_private"`
	PasteFormatLong  string `xml:"paste_format_long"`
	PasteFormatShort string `xml:"paste_format_short"`
	PasteHits        string `xml:"paste_hits"`
}

func (ps *PasteServiceOp) Create(ctx context.Context, pr *PasteRequest) (*Paste, *Response, error) {

	pr.APIDevKey = ps.client.APIDevKey
	pr.APIUserKey = ps.client.APIUserKey
	pr.APIOption = "paste"

	req, err := ps.client.NewRequest(ctx, http.MethodPost, pasteBasePath, pr)
	if err != nil {
		return nil, nil, err
	}

	resp, err := ps.client.Do(ctx, req.request)
	if err != nil {
		return nil, nil, err
	}

	return &Paste{PasteURL: resp.response.String()}, resp, nil
}

func (ps *PasteServiceOp) Delete(ctx context.Context, pdr *PasteDeleteRequest) (*Response, error) {
	pdr.APIDevKey = ps.client.APIDevKey
	pdr.APIUserKey = ps.client.APIUserKey
	pdr.APIOption = "delete"

	req, err := ps.client.NewRequest(ctx, http.MethodPost, pasteBasePath, pdr)
	if err != nil {
		return nil, err
	}

	resp, err := ps.client.Do(ctx, req.request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (ps *PasteServiceOp) List(ctx context.Context, plr *PasteListRequest) ([]Paste, *Response, error) {

	plr.APIDevKey = ps.client.APIDevKey
	plr.APIUserKey = ps.client.APIUserKey
	plr.APIOption = "list"

	req, err := ps.client.NewRequest(ctx, http.MethodPost, pasteBasePath, plr)
	if err != nil {
		return nil, nil, err
	}

	resp, err := ps.client.Do(ctx, req.request)
	if err != nil {
		return nil, nil, err
	}

	var pastes []Paste
	d := xml.NewDecoder(bytes.NewBufferString(resp.response.String()))
	for {
		var t Paste
		err := d.Decode(&t)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		pastes = append(pastes, t)
	}
	return pastes, resp, nil
}

func (ps *PasteServiceOp) GetRaw(ctx context.Context, rpr *PasteGetRawRequest) ([]byte, *Response, error) {
	rpr.APIDevKey = ps.client.APIDevKey
	rpr.APIUserKey = ps.client.APIUserKey
	rpr.APIOption = "show_paste"

	req, err := ps.client.NewRequest(ctx, http.MethodPost, rawPasteBasePath, rpr)
	if err != nil {
		return nil, nil, err
	}

	resp, err := ps.client.Do(ctx, req.request)
	if err != nil {
		return nil, nil, err
	}
	return resp.response.Body(), resp, nil
}
