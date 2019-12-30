package yandex

type Settlement struct {
	Name     string      `json:"title"`
	Codes    interface{} `json:"codes"`
	Stations []Station   `json:"stations"`
}
