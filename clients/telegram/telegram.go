package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	ctx    context.Context
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
	sendMessage method = "sendMessage"
	getMe       method = "getMe"
)

type queryString map[string]string

// https://api.telegram.org/bot<token>/METHOD_NAME
func New(ctx context.Context) (*Client, error) {
	token, err := mustToken(tokenFile)
	if err != nil {
		return nil, err
	}

	c := &Client{
		host:   host,
		path:   "bot" + token,
		offset: 0,
		ctx:    ctx,
		client: http.Client{
			Timeout: (timeout + 2) * time.Second,
		},
	}

	resp, err := c.doRequest(getMe, nil)
	if err != nil {
		return nil, fmt.Errorf("getMe error %v", err)
	}

	var result CheckStatus
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("JSON unmarshal error %v", err)
	}

	if !result.Ok {
		return nil, fmt.Errorf("invalid token")
	} else {
		log.Printf("Bot %v identified...", result.Bot.FirstName)
	}
	return c, nil
}

// TODO deprecated. Inplemented in external logic.
func (c *Client) Start() {
	log.Println("Bot started...")
	for {
		ubdates, err := c.Update()
		if c.ctx.Err() != nil {
			return
		}
		if err != nil {
			log.Println(err)
		}

		for _, upd := range ubdates {
			fmt.Println(
				time.Unix(upd.Message.Time, 0),
				upd.Message.User.UserName,
				upd.Message.Chat.Type,
				upd.Message.Text,
				// upd.Message.ID,
			)
			if err := c.SendMessage(
				upd.Message.Chat,
				upd.Message.Text,
			); err != nil {
				log.Println(err)
			}
		}
	}
}

func (c *Client) Update() ([]Update, error) {
	query := map[string]string{
		"offset":  strconv.FormatInt(c.offset, 10),
		"limit":   strconv.Itoa(messagesLimit),
		"timeout": strconv.FormatInt(timeout, 10),
	}

	resp, err := c.doRequest(getUpdates, query)
	if err != nil {
		// TODO No need to exit, need to retry
		return nil, fmt.Errorf("fail to update %v", err)
	}

	var updates UpdateResponse
	if err := json.Unmarshal(resp, &updates); err != nil {
		return nil, fmt.Errorf("JSON unmarshal error %v", err)
	}
	// TODO BY this in future we can check en internal Telegram API error
	if !updates.Ok {
		return nil, nil
	}

	for _, upd := range updates.Result {
		if upd.ID >= c.offset {
			c.offset = upd.ID + 1
		}
	}

	return updates.Result, nil
}

func (c *Client) SendMessage(to Chat, msg string) error {
	query := map[string]string{
		"chat_id": strconv.FormatInt(to.ID, 10),
		"text":    msg,
	}

	_, err := c.doRequest(sendMessage, query)
	if err != nil {
		return fmt.Errorf("fail to send message %v", err)
	}

	return nil
}

func (c *Client) Ctx() error {
	return c.ctx.Err()
}

func (c *Client) doRequest(m method, query queryString) ([]byte, error) {
	q := url.Values{}
	for i, val := range query {
		q.Add(i, val)
	}

	u := url.URL{
		Scheme: protocol,
		Host:   host,
		Path:   path.Join(c.path, string(m)),
	}

	req, err := http.NewRequestWithContext(c.ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("request creation fail %v", err)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response error %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("body read error %v", err)
	}

	return body, nil
}

func mustToken(filePath string) (string, error) {
	content, err := os.ReadFile(tokenFile)
	if err != nil {
		return "", fmt.Errorf("token file error %v", err)
	}

	return strings.TrimSpace(string(content)), nil
}
