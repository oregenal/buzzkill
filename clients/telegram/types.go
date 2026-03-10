package telegram

type Response struct {
	OK bool `json:"ok"`
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
	Text string `json:"text"`
	Time int64  `json:"date"`
}
