package yandex

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	dateFormat  = "2006-01-02"
	scheme      = "https"
	defaultHost = "api.rasp.yandex.net"
	apiVersion  = "v3.0"
)

type Client interface {
	// Расписание рейсов по станции
	Schedules(ctx context.Context, req SchedulesRequest) (*SchedulesResponse, error)
	// Список всех доступных станций
	StationsList(ctx context.Context) (*StationsListResponse, error)
	// Расписание рейсов между станциями
	Search(ctx context.Context, req SearchRequest) (*SearchResponse, error)
	// Список станций следования
	Thread(ctx context.Context, req ThreadRequest) (*ThreadResponse, error)
}

type client struct {
	client *http.Client
	cfg    *Config
}

// New return client
func New(cfg *Config) Client {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	httpClient := &http.Client{
		Timeout:   cfg.Timeout,
		Transport: transport,
	}
	return &client{
		client: httpClient,
		cfg:    cfg,
	}
}

// NewWithDefaultConfig return client with default cfg
func NewWithDefaultConfig(apiKey string) Client {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	httpClient := &http.Client{
		Timeout:   5 * time.Second,
		Transport: transport,
	}

	return &client{
		client: httpClient,
		cfg: &Config{
			Host:    defaultHost,
			ApiKey:  apiKey,
			Format:  JsonFormat,
			Lang:    Ru,
			Version: apiVersion,
		},
	}
}

type SchedulesRequest struct {
	Station       string        //
	Time          time.Time     //
	TransportType transportType //
}

type SchedulesResponse struct {
	Date              string      `json:"date"`               // Дата, на которую получен список рейсов.
	Pagination        interface{} `json:"pagination"`         // Информация о постраничном выводе найденных рейсов.
	Station           Station     `json:"station"`            // Информация об указанной в запросе станции.
	Schedule          []Schedule  `json:"schedule"`           // Список рейсов.
	ScheduleDirection interface{} `json:"schedule_direction"` // Код и название запрошенного направления рейсов.
	Directions        interface{} `json:"directions"`         // Коды и названия возможных направлений движения электричек по станции.
}

func (c *client) Schedules(ctx context.Context, req SchedulesRequest) (*SchedulesResponse, error) {
	u := url.URL{
		Scheme: scheme,
		Host:   c.cfg.Host,
		Path:   c.cfg.Version + "/schedule/",
	}

	q := u.Query()
	q.Set("apikey", c.cfg.ApiKey)
	q.Set("format", c.cfg.Format.String())
	q.Set("lang", c.cfg.Lang.String())
	q.Set("transport_type", req.TransportType.String())
	q.Set("station", req.Station)
	q.Set("date", req.Time.Format(dateFormat))

	u.RawQuery = q.Encode()

	var resp SchedulesResponse
	if err := c.get(ctx, u.String(), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

type StationsListResponse struct {
	Countries []Country `json:"countries"`
}

func (c *client) StationsList(ctx context.Context) (*StationsListResponse, error) {
	u := url.URL{
		Scheme: scheme,
		Host:   c.cfg.Host,
		Path:   c.cfg.Version + "/stations_list/",
	}

	q := u.Query()
	q.Set("apikey", c.cfg.ApiKey)
	q.Set("format", c.cfg.Format.String())
	q.Set("lang", c.cfg.Lang.String())

	u.RawQuery = q.Encode()

	var resp StationsListResponse
	if err := c.get(ctx, u.String(), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

type SearchRequest struct {
	From string
	To   string
	Date time.Time
}

type SearchResponse struct {
	Pagination       interface{}   `json:"pagination"`
	IntervalSegments []interface{} `json:"interval_segments"`
	Segments         []Segment     `json:"segments"`
	Search           interface{}   `json:"search"`
}

func (c *client) Search(ctx context.Context, req SearchRequest) (*SearchResponse, error) {
	if req.From == "" || req.To == "" || req.Date.IsZero() {
		return nil, errors.New("one of required request param are missing")
	}

	u := url.URL{
		Scheme: scheme,
		Host:   c.cfg.Host,
		Path:   c.cfg.Version + "/search/",
	}

	q := u.Query()
	q.Set("apikey", c.cfg.ApiKey)
	q.Set("format", c.cfg.Format.String())
	q.Set("lang", c.cfg.Lang.String())
	q.Set("from", req.From)
	q.Set("to", req.To)
	q.Set("date", req.Date.Format(dateFormat))

	u.RawQuery = q.Encode()

	var resp SearchResponse
	if err := c.get(ctx, u.String(), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

type ThreadRequest struct {
	UID  string
	From string
	To   string
}

type ThreadResponse struct {
	Stops []Stop `json:"stops"`
}

func (c *client) Thread(ctx context.Context, req ThreadRequest) (*ThreadResponse, error) {
	if req.UID == "" {
		return nil, errors.New("uid are missing")
	}

	u := url.URL{
		Scheme: scheme,
		Host:   c.cfg.Host,
		Path:   c.cfg.Version + "/thread/",
	}

	q := u.Query()
	q.Set("apikey", c.cfg.ApiKey)
	q.Set("format", c.cfg.Format.String())
	q.Set("lang", c.cfg.Lang.String())
	q.Set("uid", req.UID)
	q.Set("from", req.From)
	q.Set("to", req.To)

	u.RawQuery = q.Encode()

	var resp ThreadResponse
	if err := c.get(ctx, u.String(), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *client) get(ctx context.Context, url string, resp interface{}) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	httpResp, err := c.client.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(httpResp.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf("%d status code: %s", httpResp.StatusCode, string(body))
	}

	switch c.cfg.Format {
	case JsonFormat:
		return json.NewDecoder(httpResp.Body).Decode(resp)
	case XmlFormat:
		return xml.NewDecoder(httpResp.Body).Decode(resp)
	default:
		return errors.New("format unsupported")
	}
}