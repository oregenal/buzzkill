package telegram

type Response struct {
	Ok     bool   `json:"ok"`
	Result Update `json:"result"`
}

type UpdateResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID      int64   `json:"update_id"`
	Message Message `json:"message"`
}

type Message struct {
	ID   int64  `json:"message_id"`
	User User   `json:"from"`
	Text string `json:"text"`
	Time int64  `json:"date"`
	Chat Chat   `json:"chat"`
}

type User struct {
	ID        int64  `json:"id"`
	Bot       bool   `json:"is_bot"`
	FisrtName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"username"`
	Lang      string `json:"language_code"`
}

type Chat struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	UserName  string `json:"username"`
	FisrtName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
