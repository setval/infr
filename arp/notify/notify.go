package notify

const (
	NotifySourceDiscord = "discord"
	NotifySourceVK      = "vk"
)

type Notify struct {
	Source string
	Data   Data
}

type Data struct {
	Recipient string
	Text      string
}
