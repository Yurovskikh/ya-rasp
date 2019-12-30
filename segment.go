package yandex

import "time"

type Segment struct {
	Arrival           time.Time   `json:"arrival"`
	From              Station     `json:"from"`
	Thread            Thread      `json:"thread"`
	DeparturePlatform string      `json:"departure_platform"`
	Departure         time.Time   `json:"departure"`
	Stops             string      `json:"stops"`
	DepartureTerminal interface{} `json:"departure_terminal"`
	To                Station     `json:"to"`
	HasTransfers      bool        `json:"has_transfers"`
	TicketsInfo       interface{} `json:"tickets_info"`
	Duration          float64     `json:"duration"`
	ArrivalTerminal   string      `json:"arrival_terminal"`
	StartDate         string      `json:"start_date"`
	ArrivalPlatform   string      `json:"arrival_platform"`
}
