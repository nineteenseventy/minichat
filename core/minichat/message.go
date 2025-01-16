package minichat

type MessageAttachment struct {
	Id        string  `json:"id"`
	MessageId string  `json:"messageId"`
	Type      string  `json:"type"`
	Filename  string  `json:"filename"`
	Url       *string `json:"url"`
}

type Message struct {
	Id          string              `json:"id"`
	ChannelId   string              `json:"channelId"`
	AuthorId    string              `json:"authorId"`
	Content     string              `json:"content"`
	Timestamp   *string             `json:"timestamp"`
	Read        bool                `json:"read"`
	Attachments []MessageAttachment `json:"attachments"`
}

type BaseMessageAttachment struct {
	Type     string `json:"type"`
	Filename string `json:"filename"`
}

type MessageBase struct {
	Content     string                  `json:"content"`
	Attachments []BaseMessageAttachment `json:"attachments"`
}
