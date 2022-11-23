package ulog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	FormatRaw  = "raw"
	FormatJson = "json"
)

type Client struct {
	settings Settings
	client   *http.Client
}

type Settings struct {
	Host     string
	Token    string
	Login    string
	Password string
}

func New(settings Settings) *Client {
	return &Client{
		settings: settings,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type Parameters struct {
	ServerID int
	Date     time.Time
	UserID   int
	Type     string
	Format   string
	Listen   int
	IsVC     bool
}

type Log struct {
	Date time.Time `json:"date"`
	Text string    `json:"text"`
	Raw  string    `json:"raw"`
}

func (c *Client) SetAuth(login, password string) {
	c.settings.Login = login
	c.settings.Password = password
}

func (c *Client) DownloadRaw(parameters Parameters) (string, error) {
	parameters.Format = FormatRaw
	body, err := c.doRequest(parameters)
	if err != nil {
		return "", err
	}
	defer body.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(body); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (c *Client) DownloadLogs(parameters Parameters) ([]Log, error) {
	parameters.Format = FormatJson
	body, err := c.doRequest(parameters)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var logs []Log
	if err := json.NewDecoder(body).Decode(&logs); err != nil {
		return nil, err
	}

	return logs, nil
}

func (c *Client) doRequest(parameters Parameters) (io.ReadCloser, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(map[string]interface{}{
		"token":     c.settings.Token,
		"login":     c.settings.Login,
		"password":  c.settings.Password,
		"user_id":   parameters.UserID,
		"server_id": parameters.ServerID,
		"date":      parameters.Date,
		"type":      parameters.Type,
		"format":    parameters.Format,
		"listen":    parameters.Listen,
		"is_vc":     parameters.IsVC,
	}); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.settings.Host+"/api/ulog/parse", &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code is %d", resp.StatusCode)
	}

	return resp.Body, nil
}
