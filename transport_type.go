package yandex

type transportType string

const (
	Plane      transportType = "plane"
	Train      transportType = "train"
	Suburban   transportType = "suburban"
	Bus        transportType = "bus"
	Water      transportType = "water"
	Helicopter transportType = "helicopter"
)

func (t transportType) String() string {
	return string(t)
}
