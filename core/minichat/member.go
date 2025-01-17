package minichat

type Member struct {
	ID       string  `json:"id"`
	UserId   string  `json:"userId"`
	Username string  `json:"username"`
	Picture  *string `json:"picture"`
}
