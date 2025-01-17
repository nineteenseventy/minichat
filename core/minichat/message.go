package minichat

type MessageAttachment struct {
	Id        string  `json:"id"`
	MessageId string  `json:"messageId"`
	Type      string  `json:"type"`
	Filename  string  `json:"filename"`
	Url       *string `json:"url"`
}

type Mention struct {
	UserId   string `json:"userId"`
	Username string `json:"username"`
}

type Message struct {
	Id          string              `json:"id"`
	ChannelId   string              `json:"channelId"`
	AuthorId    string              `json:"authorId"`
	Content     string              `json:"content"`
	Timestamp   *string             `json:"timestamp"`
	Read        bool                `json:"read"`
	Mentions    []Mention           `json:"mentions"`
	Attachments []MessageAttachment `json:"attachments"`
}

type BaseMessageAttachment struct {
	Type     string `json:"type"`
	Filename string `json:"filename"`
}

type BaseMessage struct {
	Content          string                  `json:"content"`
	MentionedUserIds []string                `json:"mentionedUserIds"`
	Attachments      []BaseMessageAttachment `json:"attachments"`
}
