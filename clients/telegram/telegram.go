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
)

type method string
const (
	getUpdates = "getUpdates"
	sendMessage = "sendMessage"
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
	query := make(map[string]string, 2)
	query["offset"] = strconv.Itoa(0)
	query["limit"] = strconv.Itoa(100)
	resp, err := c.doRequest(getUpdates, query)
	if err != nil {
		return nil, fmt.Errorf("Fail to update: %v", err)
	}

	var res UpdateResponse
	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("JSON Unmarshal error: %v", err)
	}
	// fmt.Println(res)
	if !res.Ok {
		return nil, nil
	}

	return res.Result, nil
}

func (c *Client) SendMessage(msg string) error {
	// "https://api.telegram.org/bot1380900708:AAHvBoVcgUSQwwdnRd0QPgSV48QHb-oHW2M/sendMessage?chat_id=220630034&text=Hello!"
	query := make(map[string]string, 2)
	query["chat_id"] = strconv.Itoa(220630034)
	query["text"] = msg
	_, err := c.doRequest(sendMessage, query)
	if err != nil {
		return fmt.Errorf("Fail to send message: %v", err)
	}
	return nil
}

func (c *Client) doRequest(m method, query map[string]string) ([]byte, error) {
	q := url.Values{}
	for i, val := range query {
		q.Add(i, val)
	}
	// fmt.Println(q)
	u := url.URL{
		Scheme: "https",
		Host:   host,
		Path:   path.Join(c.path, string(m)),
	}
	// fmt.Println(u)
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Request creation fail %v", err)
	}
	req.URL.RawQuery = q.Encode()
	fmt.Println(req)
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
	return body, err
}

func mustToken(filePath string) (string, error) {
	content, err := os.ReadFile(tokenFile)
	if err != nil {
		return "", fmt.Errorf("Token file error %v", err)
	}
	return string(content[:len(content)-1]), nil
}
