package loader

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
)

const (
	FormatRaw  = "raw"
	FormatJson = "json"
)

const (
	EngineUlog      = "ulog"
	EngineGamePanel = "gamepanel"
)

var ErrUnknownEngine = errors.New("unknown engine")

func (c *Client) downloadRaw(r io.ReadCloser) (string, error) {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(r); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (c *Client) downloadLogJSON(r io.ReadCloser) ([]Log, error) {
	var logs []Log
	if err := json.NewDecoder(r).Decode(&logs); err != nil {
		return nil, err
	}
	return logs, nil
}

func (c *Client) downloadUserJSON(r io.ReadCloser) ([]User, error) {
	var users []User
	if err := json.NewDecoder(r).Decode(&users); err != nil {
		return nil, err
	}
	return users, nil
}

func (c *Client) loadUlog(parameters UlogParameters) (io.ReadCloser, error) {
	params := c.createParams(parameters.Parameters)
	params = c.addUlogParams(params, parameters)

	body, err := c.doRequest(EngineUlog, params)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *Client) loadGamePanel(parameters GamePanelParameters) (io.ReadCloser, error) {
	params := c.createParams(parameters.Parameters)
	params = c.addGamePanelParams(params, parameters)

	body, err := c.doRequest(EngineGamePanel, params)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *Client) download(parameters interface{}) (io.ReadCloser, error) {
	var (
		body io.ReadCloser
		err  error
	)

	switch p := parameters.(type) {
	case UlogParameters:
		body, err = c.loadUlog(p)
	case GamePanelParameters:
		body, err = c.loadGamePanel(p)
	default:
		return nil, ErrUnknownEngine
	}

	return body, err
}

func (c *Client) DownloadRaw(parameters interface{}) (string, error) {
	body, err := c.download(parameters)
	if err != nil {
		return "", err
	}
	defer body.Close()
	return c.downloadRaw(body)
}

func (c *Client) DownloadLogJSON(parameters interface{}) ([]Log, error) {
	body, err := c.download(parameters)
	if err != nil {
		return nil, err
	}
	defer body.Close()
	return c.downloadLogJSON(body)
}

func (c *Client) DownloadUserJSON(parameters interface{}) ([]User, error) {
	body, err := c.download(parameters)
	if err != nil {
		return nil, err
	}
	defer body.Close()
	return c.downloadUserJSON(body)
}
