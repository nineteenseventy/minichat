package minichat

type BaseGroupChannel struct {
	Title string `json:"title"`
}

type GroupChannel struct {
	BaseChannel
	BaseGroupChannel
	Users []User `json:"users"`
}

type NewGroupChannel struct {
	BaseGroupChannel
	UserIds []string `json:"userIds"`
}

type UpdateGroupChannel struct {
	BaseGroupChannel
}
