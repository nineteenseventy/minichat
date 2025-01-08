package minichat

type BaseChannel struct {
	Id          *string `json:"id"`
	Type        string  `json:"type"`
	CreatedAt   *string `json:"createdAt"`
	UnreadCount *int    `json:"unreadCount"`
}

type ChannelPrivate struct {
	BaseChannel
	Title string `json:"title"`
}

type ChannelPublic struct {
	BaseChannel
	Title       string  `json:"title"`
	Description *string `json:"description"`
}

type Channel struct {
	Id           string    `json:"id"`
	Type         string    `json:"type"`
	Title        string    `json:"title"`
	Description  *string   `json:"description"`
	CreatedAt    *string   `json:"createdAt"`
	LastMessages []Message `json:"lastMessages"`
}
