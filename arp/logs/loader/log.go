package loader

import "time"

type Log struct {
	Date time.Time `json:"date"`
	Text string    `json:"text"`
	Raw  string    `json:"raw"`
}
