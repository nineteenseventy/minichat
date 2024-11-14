package minichat

type User struct {
	ID       string  `json:"id"`
	Username string  `json:"username"`
	Picture  *string `json:"profilePicture"`
}

type UserProfile struct {
	ID       string  `json:"id"`
	Username string  `json:"username"`
	Picture  *string `json:"picture"`
	Bio      *string `json:"bio"`
	Color    *string `json:"color"`
}
