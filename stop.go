package yandex

type Stop struct {
	Arrival   string  `json:"arrival"`
	Departure string  `json:"departure"`
	Terminal  string  `json:"terminal"`
	Station   Station `json:"station"`
	StopTime  int     `json:"stop_time"`
	Duration  float64 `json:"duration"`
}
