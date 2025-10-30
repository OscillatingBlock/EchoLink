// pkg/ml/client.go
package ml

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Client struct {
	baseURL string
	client  *http.Client
}

func NewClient(url string) *Client {
	return &Client{baseURL: url, client: &http.Client{}}
}

type Request struct {
	BotID     string `json:"bot_id"`
	UserInput string `json:"user_input"`
	Goal      string `json:"goal"`
	Context   string `json:"context"`
	Webhook   string `json:"webhook"`
}

type Response struct {
	Response string `json:"response"`
	EndCall  bool   `json:"end_call"`
}

func (c *Client) Process(req Request) (Response, error) {
	var res Response
	body, _ := json.Marshal(req)
	httpRes, err := c.client.Post(c.baseURL+"/ml/process", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return res, err
	}
	defer httpRes.Body.Close()
	json.NewDecoder(httpRes.Body).Decode(&res)
	return res, nil
}
