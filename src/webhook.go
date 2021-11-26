package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
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
	Guild_id       string `json:"guild_id"`
	Chanel_id      string `json:"channel_id"`
	UserData       User   `json:"user"`
	Name           string `json:"name"`
	Avatar         string `json:"avatar"`
	Token          string `json:"token"`
	Application_id string `json:"application_id"`
	Source_guild   string `json:"source_guild"`
	Source_channel string `json:"source_channel"`
	URL            string `json:"url"`
}

type Icon struct {
	IconURL      string `json:"icon_url"`
	ProxyIconURL string `json:"proxy_icon_url"`
}

type EmbedFooter struct {
	Text string `json:"text"`
	Icon
}

type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type EmbedProvider struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type EmbedAuthor struct {
	EmbedProvider
	Icon
}

// Video, thumbnail and image have the same JSON schema
type EmbedImage struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
}
type EmbedVideo EmbedImage
type EmbedThumbnail EmbedImage

// TODO
type Embed struct {
	Title       string         `json:"title"`
	Type        string         `json:"type"`
	Description string         `json:"description"`
	URL         string         `json:"url"`
	Timestamp   time.Time      `json:"timestamp"`
	Color       int            `json:"color"`
	Footer      EmbedFooter    `json:"footer"`
	Image       EmbedImage     `json:"image"`
	Thumbnail   EmbedThumbnail `json:"thumbnail"`
	Video       EmbedVideo     `json:"video"`
	Provider    EmbedProvider  `json:"provider"`
	Author      EmbedAuthor    `json:"author"`
	Fields      []EmbedField   `json:"fields"`
}

func (wh *Webhook) Connect(url string) {
	wh.URL = url
	res, err := http.Get(url)
	if err != nil {
		fmt.Print("Error connecting to webhook")
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	fmt.Printf("Body: %s\n", string(body))

	json.Unmarshal(body, wh)
}

func (wh *Webhook) SendRaw(data string) {
	apiUrl := fmt.Sprintf("https://discord.com/api/webhooks/%s/%s", wh.Id, wh.Token)
	req, _ := http.NewRequest("POST", apiUrl, bytes.NewBuffer([]byte(data)))
	req.Header.Set("Content-Type", "application/json")
	_, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Printf("Websocket request error: %s\n", err)
		return
	}
}

func (wh *Webhook) Send(message string) {
	msgJson := fmt.Sprintf(`{
		"content": "%s"
	}`, message)

	wh.SendRaw(msgJson)
}

func (wh *Webhook) SendEmbed(embed []Embed) {
	data, _ := json.Marshal(embed)
	msgJson := fmt.Sprintf(`{
		"embeds": [%s]
	}`, string(data))

	wh.SendRaw(msgJson)
}
