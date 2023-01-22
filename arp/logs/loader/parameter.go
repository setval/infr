package loader

import "time"

type Parameters struct {
	Engine   string    `json:"engine"`
	Type     string    `json:"type"`
	Token    string    `json:"token"`
	Format   string    `json:"format"`
	Listen   int       `json:"listen"`
	ServerID int       `json:"server_id"`
	Full     bool      `json:"full"`
	Date     time.Time `json:"date"`
	Login    string    `json:"login"`
	Password string    `json:"password"`
}

type UlogParameters struct {
	Parameters
	Text    string `json:"text"`
	UserID  int    `json:"user_id"`
	IsVC    bool   `json:"is_vc"`
	IsAdmin bool   `json:"is_admin"`
}

type GamePanelParameters struct {
	Parameters
	Host     string   `json:"host"`
	Secret   string   `json:"secret"`
	Nickname string   `json:"nickname"`
	Reason   string   `json:"reason"`
	Params   []string `json:"params"`
}
