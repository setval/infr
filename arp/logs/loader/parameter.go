package loader

import "time"

type Parameters struct {
	Engine   string    `json:"engine"`
	Type     string    `json:"type"`
	Token    string    `json:"token"`
	Format   string    `json:"format"`
	Listen   int       `json:"listen"`
	ServerID int       `json:"server_id"`
	Date     time.Time `json:"date"`
	Login    string    `json:"login"`
	Password string    `json:"password"`
}

type UlogParameters struct {
	Parameters
	UserID  int  `json:"user_id"`
	IsVC    bool `json:"is_vc"`
	IsAdmin bool `json:"is_admin"`
}

type GamePanelParameters struct {
	Parameters
	Secret   string `json:"secret"`
	Nickname string `json:"nickname"`
}
