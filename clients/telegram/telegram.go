package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
)

const (
	tokenFile  = ".tgToken"
	host       = "api.telegram.org"
	getUpdates = "getUpdates"
)

type Client struct {
	host   string
	path   string
	client http.Client
}

// https://api.telegram.org/bot<token>/METHOD_NAME

func New() (*Client, error) {
	token, err := mustToken(tokenFile)
	if err != nil {
		return nil, err
	}
	return &Client{
		host:   host,
		path:   "bot" + token,
		client: http.Client{},
	}, nil
}

func (c *Client) Update() ([]Update, error) {
	v := url.Values{}
	// v.Add("offset", strconv.Itoa(1))
	v.Add("limit", strconv.Itoa(10))
	// fmt.Println(q)
	u := url.URL{
		Scheme: "https",
		Host:   host,
		Path:   path.Join(c.path, getUpdates),
	}
	// fmt.Println(u)
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Request creation fail %v", err)
	}
	req.URL.RawQuery = v.Encode()
	// fmt.Println(req)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Response error %v", err)
	}
	// fmt.Println(resp)
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Body read error: %v", err)
	}
	// fmt.Println(string(body))

	var res UpdateResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("JSON Unmarshal error: %v", err)
	}
	// fmt.Println(res)
	if !res.Ok {
		return nil, nil
	}

	return res.Result, nil
}

func (c *Client) SendMessage() {
}

func mustToken(filePath string) (string, error) {
	content, err := os.ReadFile(tokenFile)
	if err != nil {
		return "", fmt.Errorf("Token file error %v", err)
	}
	return string(content[:len(content)-1]), nil
}
