package yandex

type format string

func (f format) String() string {
	return string(f)
}

const (
	JsonFormat format = "json"
	XmlFormat  format = "xml"
)
