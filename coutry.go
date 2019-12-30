package yandex

type Country struct {
	Regions []Region `json:"regions"`
	Name    string   `json:"title"`
}
