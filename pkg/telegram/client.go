package telegram

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type Message struct {
	Sender  string
	Title   string
	Message string
}

type Client struct {
	apiKey    string
	channelId string
	sender    string
}

func (c *Client) SendMessage(message Message) error {

	if message.Sender == "" {
		hostname, err := os.Hostname()
		if err == nil {
			message.Sender = hostname
		}
	}
	text := fmt.Sprintf("*From:%s*\n\n*%s*\n%s", message.Sender, message.Title, message.Message)
	query := fmt.Sprintf("chat_id=%s&text=%s&parse_mode=MarkdownV2", c.channelId, url.QueryEscape(text))
	fullUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?%s", c.apiKey, query)
	_, err := http.Get(fullUrl)
	return err
}

func NewClient(apiKey, channelId, sender string) *Client {
	cl := new(Client)
	cl.apiKey = apiKey
	cl.channelId = channelId
	cl.sender = sender
	return cl
}
