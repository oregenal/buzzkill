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
	"strings"
	"time"
)

type Client struct {
	host   string
	path   string
	client http.Client
	offset int64
}

const (
	tokenFile     = ".tgToken"
	host          = "api.telegram.org"
	protocol      = "https"
	timeout       = 30
	messagesLimit = 100
)

type method string

const (
	getUpdates  method = "getUpdates"
	sendMessage        = "sendMessage"
)

type queryString map[string]string

// https://api.telegram.org/bot<token>/METHOD_NAME

func New() (*Client, error) {
	token, err := mustToken(tokenFile)
	if err != nil {
		return nil, err
	}

	return &Client{
		host: host,
		path: "bot" + token,
		client: http.Client{
			Timeout: (timeout + 2) * time.Second,
		},
		offset: 0,
	}, nil
}

func (c *Client) Update() ([]Update, error) {
	query := make(map[string]string, 2)
	query["offset"] = strconv.FormatInt(c.offset, 10)
	query["limit"] = strconv.Itoa(messagesLimit)
	query["timeout"] = strconv.FormatInt(timeout, 10)

	resp, err := c.doRequest(getUpdates, query)
	if err != nil {
		// TODO No need to exit, need to retry
		return nil, fmt.Errorf("Fail to update: %v", err)
	}

	var updates UpdateResponse
	if err := json.Unmarshal(resp, &updates); err != nil {
		return nil, fmt.Errorf("JSON Unmarshal error: %v", err)
	}
	// fmt.Println(updates)
	// TODO BY this in future we can check en internal Telegram API error
	if !updates.Ok {
		return nil, nil
	}

	for _, upd := range updates.Result {
		if upd.ID >= c.offset {
			c.offset = upd.ID + 1
			// fmt.Println(upd.ID)
		}
	}

	return updates.Result, nil
}

func (c *Client) SendMessage(to Chat, msg string) error {
	query := make(map[string]string, 2)
	query["chat_id"] = strconv.FormatInt(to.ID, 10)
	query["text"] = msg

	_, err := c.doRequest(sendMessage, query)
	if err != nil {
		return fmt.Errorf("Fail to send message: %v", err)
	}

	return nil
}

func (c *Client) doRequest(m method, query queryString) ([]byte, error) {
	q := url.Values{}
	for i, val := range query {
		q.Add(i, val)
	}
	// fmt.Println(q)

	u := url.URL{
		Scheme: protocol,
		Host:   host,
		Path:   path.Join(c.path, string(m)),
	}
	// fmt.Println(u)

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Request creation fail %v", err)
	}
	req.URL.RawQuery = q.Encode()
	// fmt.Println(req)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Response error %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	// fmt.Println(resp)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Body read error: %v", err)
	}
	// fmt.Println(string(body))

	return body, nil
}

func mustToken(filePath string) (string, error) {
	content, err := os.ReadFile(tokenFile)
	if err != nil {
		return "", fmt.Errorf("Token file error %v", err)
	}

	return strings.TrimSpace(string(content)), nil
}
