package minichat

type User struct {
	ID       string  `json:"id"`
	Username string  `json:"username"`
	Picture  *string `json:"picture"`
}

type UserProfile struct {
	ID       string  `json:"id"`
	Username string  `json:"username"`
	Picture  *string `json:"picture"`
	Bio      *string `json:"bio"`
}

type PatchUserProfile struct {
	Username *string `json:"username"`
	Picture  *string `json:"picture"`
	Bio      *string `json:"bio"`
}

type UserStatus struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}
