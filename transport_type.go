package yandex

type TransportType string

const (
	Plane      TransportType = "plane"
	Train      TransportType = "train"
	Suburban   TransportType = "suburban"
	Bus        TransportType = "bus"
	Water      TransportType = "water"
	Helicopter TransportType = "helicopter"
)

func (t TransportType) String() string {
	return string(t)
}
