package yandex

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestClient_Search(t *testing.T) {
	t.Skip()
	client := NewWithDefaultConfig("") // set api key
	resp, err := client.Search(context.TODO(), SearchRequest{
		From: "s2006004",
		To:   "s9602494",
		Date: time.Now().AddDate(0, 0, 5),
	})
	if err != nil {
		t.Error(err)
	}

	if len(resp.Segments) == 0 {
		t.Error("empty segments")
	}
}

func TestClient_Schedules(t *testing.T) {
	t.Skip()
	client := NewWithDefaultConfig("") // set api key
	reps, err := client.Schedules(context.TODO(), SchedulesRequest{
		Station:       "s2006004",
		Time:          time.Now(),
		TransportType: Train,
	})
	if err != nil {
		t.Error(err)
	}

	if len(reps.Schedule) == 0 {
		t.Error("empty schedules")
	}
}

func TestClient_StationsList(t *testing.T) {
	t.Skip()
	client := NewWithDefaultConfig("") // set api key
	resp, err := client.StationsList(context.TODO())
	if err != nil {
		t.Error(err)
	}

	if len(resp.Countries) == 0 {
		t.Error("empty countries")
	}
}

func TestClient_Thread(t *testing.T) {
	t.Skip()
	client := NewWithDefaultConfig("") // set api key
	resp, err := client.Thread(context.TODO(), ThreadRequest{
		UID: "726CH_1_2",
	})
	if err != nil {
		fmt.Println(err)
		t.Fatal()
	}

	if len(resp.Stops) == 0 {
		t.Error("empty stops")
	}
}
