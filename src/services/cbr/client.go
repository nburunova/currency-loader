package cbr

import (
	"bytes"
	"context"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"golang.org/x/net/html/charset"
)

const (
	cbrURL = "http://www.cbr.ru/scripts/XML_daily.asp"
)

var (
	ErrCreateRequest = errors.New("Cannot create request for valutes")
	ErrServerError   = errors.New("Server error")
	ErrStatusNotOK   = errors.New("Status code not OK")
	ErrDoneByContext = errors.New("Done by context")
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	client  HTTPClient
	baseURL url.URL
}

func NewClient(client HTTPClient) (*Client, error) {
	if client == nil {
		client = http.DefaultClient
	}
	url, err := url.Parse(cbrURL)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot parse currency source URL %v", cbrURL)
	}
	return &Client{
		client:  client,
		baseURL: *url,
	}, nil
}

func (c *Client) Valutes(ctx context.Context) (Valutes, error) {
	req, err := http.NewRequest(http.MethodGet, c.baseURL.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create request")
	}
	var res rawValCurs
	err = c.do(ctx, req, &res)
	if err != nil {
		return nil, errors.Wrap(err, "cannot request valutes")
	}
	return res.Valutes()
}

func (c *Client) do(ctx context.Context, req *http.Request, holder interface{}) error {
	response, err := c.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "cannot do request")
	}
	if response.StatusCode != http.StatusOK {
		if response.StatusCode >= http.StatusInternalServerError {
			return errors.Wrapf(ErrServerError, "Server response %v", response.Status)
		}
		return errors.Wrapf(ErrStatusNotOK, "Server response %v", response.Status)
	}

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.Wrapf(err, "cannot read content from response for %v", req.URL)
	}
	err = response.Body.Close()
	if err != nil {
		return errors.Wrapf(err, "cannot close response to %v", req.URL)
	}
	decoder := xml.NewDecoder(bytes.NewReader(content))
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(holder)
	if err != nil {
		return errors.Wrapf(err, "error when parsing pesponse %v", string(content))
	}
	return nil
}
