package minichat

type Channel struct {
	Id        *string `json:"id"`
	Type      string  `json:"type"`
	CreatedAt string  `json:"createdAt"`
}

type ChannelPrivate struct {
	Channel
	Title string `json:"title"`
}

type ChannelPublic struct {
	Channel
	Title       string  `json:"title"`
	Description *string `json:"description"`
}
