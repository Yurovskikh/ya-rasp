package yandex

type Thread struct {
	UID              string      `json:"uid"`
	Title            string      `json:"title"`
	Number           string      `json:"number"`
	ShortTitle       string      `json:"short_title"`
	ThreadMethodLink string      `json:"thread_method_link"`
	Carrier          Carrier `json:"carrier"`
	Address          string      `json:"address"`
	Logo             string      `json:"logo"`
	Email            string      `json:"email"`
}
