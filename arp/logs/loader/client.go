package loader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
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
	Secret   string
}

func New(settings Settings) *Client {
	return &Client{
		settings: settings,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (c *Client) SetAuth(login, password string) {
	c.settings.Login = login
	c.settings.Password = password
}

func (c *Client) createParams(parameters Parameters) map[string]interface{} {
	return map[string]interface{}{
		"token":     c.settings.Token,
		"login":     c.settings.Login,
		"password":  c.settings.Password,
		"server_id": parameters.ServerID,
		"date":      parameters.Date,
		"type":      parameters.Type,
		"format":    parameters.Format,
		"listen":    parameters.Listen,
		"full":      parameters.Full,
	}
}

func (c *Client) addUlogParams(parameters map[string]interface{}, ulogParameters UlogParameters) map[string]interface{} {
	parameters["user_id"] = ulogParameters.UserID
	parameters["is_vc"] = ulogParameters.IsVC
	parameters["is_admin"] = ulogParameters.IsAdmin
	parameters["text"] = ulogParameters.Text
	return parameters
}

func (c *Client) addGamePanelParams(parameters map[string]interface{}, gamePanelParameters GamePanelParameters) map[string]interface{} {
	parameters["secret"] = c.settings.Secret
	parameters["nickname"] = gamePanelParameters.Nickname
	return parameters
}

func (c *Client) doRequest(engine string, parameters map[string]interface{}) (io.ReadCloser, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(parameters); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.settings.Host+"/api/"+engine+"/parse", &buf)
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
