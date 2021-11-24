package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type User struct {
	Id            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
	Bot           bool   `json:"bot"`
	System        bool   `json:"system"`
	Mfa_enabled   bool   `json:"mfa_enabled"`
	Banner        string `json:"banner"`
	Accent_color  int    `json:"accent_color"`
	Locale        string `json:"locale"`
	Verified      bool   `json:"verified"`
	Email         string `json:"email"`
	Flags         int    `json:"flags"`
	Premium_type  int    `json:"premium_type"`
	Public_flags  int    `json:"public_flags"`
}

type Webhook struct {
	Id             string `json:"id"`
	Wh_type        int    `json:"type"`
	Guild_id       string `json:"type"`
	Chanel_id      string `json:"channel_id"`
	UserData       User   `json:"user"`
	Name           string `json:"name"`
	Avatar         string `json:"avatar"`
	Token          string `json:"token"`
	Application_id string `json:"application_id"`
	Source_guild   string `json:"source_guild"`
	Source_channel string `json:"source_channel"`
	Url            string `json:"url"`
}

// TODO
type Embed struct {
}

func (wh *Webhook) Connect(url string) {
	wh.Url = url
	res, err := http.Get(url)
	defer res.Body.Close()
	if err != nil {
		fmt.Print("Error connecting to webhook")
	}

	body, _ := ioutil.ReadAll(res.Body)
	fmt.Printf("Body: %s\n", string(body))

	json.Unmarshal(body, wh)
}

func (wh *Webhook) SendRaw(data string) {
	api_url := fmt.Sprintf("https://discord.com/api/webhooks/%s/%s", wh.Id, wh.Token)
	req, _ := http.NewRequest("POST", api_url, bytes.NewBuffer([]byte(data)))
	req.Header.Set("Content-Type", "application/json")
	_, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Printf("Websocket request error: %s\n", err)
		return
	}
}

func (wh *Webhook) Send(message string) {
	msg_json := fmt.Sprintf(`{
		"content": "%s"
	}`, message)

	wh.SendRaw(msg_json)
}

func (wh *Webhook) SendEmbed(embed []Embed) {
	data, _ := json.Marshal(embed)
	msg_json := fmt.Sprintf(`{
		"embeds": [%s]
	}`, string(data))

	wh.SendRaw(msg_json)
}
