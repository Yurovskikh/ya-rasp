package yandex

type lang string

func (l lang) String() string {
	return string(l)
}

const (
	Ru lang = "ru_RU" // Russian
	Ua lang = "uk_UA" // Ukraine
)
