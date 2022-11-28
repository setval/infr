package content

import "encoding/json"

const (
	TypeNotify = "notify"
)

type Content struct {
	Type string
	Data json.RawMessage
}
