package minichat

type Channel struct {
	Id   *string `json:"id"`
	Type string  `json:"type"`
}

type ChannelDirect struct {
	Channel
	User1Id string `json:"user1Id"`
	User2Id string `json:"user2Id"`
}

type ChannelPrivate struct {
	Channel
	Title     string `json:"title"`
	CreatedAt string `json:"createdAt"`
}

type ChannelPublic struct {
	Channel
	Title       string  `json:"title"`
	Description *string `json:"description"`
	CreatedAt   *string `json:"createdAt"`
}
