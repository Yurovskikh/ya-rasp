package yandex

type Station struct {
	Direction     string            `json:"direction"`
	Codes         map[string]string `json:"codes"`
	Type          string            `json:"station_type"`
	Title         string            `json:"title"`
	Lng           interface{}       `json:"longitude"`
	Lat           interface{}       `json:"latitude"`
	TransportType string            `json:"transport_type"`
	Code 	      string 		`json:"code"`

	Region string
	City   string
}

func (s *Station) ExternalID() (string, bool) {
	code, ok := s.Codes["yandex_code"]
	return code, ok
}
