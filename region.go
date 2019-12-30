package yandex

// Регион страны
type Region struct {
	Settlements []Settlement `json:"settlements"`
	Name        string       `json:"title"`
}
