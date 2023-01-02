package loader

import (
	"bytes"
	"encoding/json"
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

func (c *Client) downloadRaw(r io.ReadCloser) (string, error) {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(r); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (c *Client) downloadJSON(r io.ReadCloser) ([]Log, error) {
	var logs []Log
	if err := json.NewDecoder(r).Decode(&logs); err != nil {
		return nil, err
	}
	return logs, nil
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

	body, err := c.doRequest(EngineUlog, params)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *Client) DownloadUlogRaw(parameters UlogParameters) (string, error) {
	body, err := c.loadUlog(parameters)
	if err != nil {
		return "", err
	}
	defer body.Close()
	return c.downloadRaw(body)
}

func (c *Client) DownloadUlogJSON(parameters UlogParameters) ([]Log, error) {
	body, err := c.loadUlog(parameters)
	if err != nil {
		return nil, err
	}
	defer body.Close()
	return c.downloadJSON(body)
}

func (c *Client) DownloadGamePanelRaw(parameters GamePanelParameters) (string, error) {
	body, err := c.loadGamePanel(parameters)
	if err != nil {
		return "", err
	}
	defer body.Close()
	return c.downloadRaw(body)
}

func (c *Client) DownloadGamePanelJSON(parameters GamePanelParameters) ([]Log, error) {
	body, err := c.loadGamePanel(parameters)
	if err != nil {
		return nil, err
	}
	defer body.Close()
	return c.downloadJSON(body)
}
