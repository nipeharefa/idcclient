package idcclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type (
	client struct {
		location   string
		accessKey  string
		httpClient *http.Client
		baseURL    string
	}

	IClient interface {
		Do(ctx context.Context, req *http.Request) (io.Reader, error)
	}
)

type (
	Location string
)

var (
	Jaksel Location = "jkt01"
	Jakut  Location = "jkt02"
)

func NewClient(loc Location, accessKey string) *client {
	baseURL := fmt.Sprintf("https://api.idcloudhost.com/v1/%s", loc)
	return &client{
		location:   loc.String(),
		accessKey:  accessKey,
		httpClient: &http.Client{},
		baseURL:    baseURL,
	}
}

func (c *client) Do(ctx context.Context, req *http.Request) (io.Reader, error) {

	finalURL := fmt.Sprintf("%s/%s", c.baseURL, req.URL.Path)

	req.URL, _ = url.Parse(finalURL)
	buf := new(bytes.Buffer)
	req.Header.Set("apiKey", c.accessKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (l Location) String() string {
	return string(l)
}
