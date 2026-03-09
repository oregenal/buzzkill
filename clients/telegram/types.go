package telegram

type UpdateResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID      uint    `json:"update_id"`
	Message Message `json:"message"`
}

type Message struct {
	Text string `json:"text"`
	Time uint   `json:"date"`
}
