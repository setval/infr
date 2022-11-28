package notify

import "encoding/json"

type Notify struct {
	Source string
	Data   Data
}

type Data struct {
	Recipient string
	Text      string
}

func (n *Notify) MarshalJSON() ([]byte, error) {
	return json.Marshal(n)
}

func (n *Notify) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, n)
}
