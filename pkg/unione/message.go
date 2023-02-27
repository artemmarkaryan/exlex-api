package unione

type Recipient struct {
	Email string `json:"email"`
}

type Body struct {
	Html      string `json:"html"`
	Plaintext string `json:"plaintext"`
}

type Message struct {
	FromEmail  string      `json:"from_email"`
	FromName   string      `json:"from_name"`
	Recipients []Recipient `json:"recipients"`
	Body       Body        `json:"body"`
	Subject    string      `json:"subject"`
}

type request struct {
	Message `json:"message"`
}
