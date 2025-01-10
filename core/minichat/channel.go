package minichat

type Channel struct {
	Id          *string `json:"id"`
	Type        string  `json:"type"`
	CreatedAt   *string `json:"createdAt"`
	UnreadCount *int    `json:"unreadCount"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
}

type NewGroupChannel struct {
	Title   string   `json:"title"`
	UserIds []string `json:"userIds"`
}
