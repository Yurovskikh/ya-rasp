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
	"strconv"
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
	// Список ближайших станций
	NearestStations(ctx context.Context, req NearestStationsRequest) (*NearestStationsResponse, error)
	//
	NearestCity(ctx context.Context, req NearestCityRequest) (*NearestCityResponse, error)
}

type client struct {
	client *http.Client
	cfg    *Config
	keys   []string
	idxKey int
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
	return NewWithPoolKey(apiKey)
}

func NewWithPoolKey(keys ...string) Client {
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
		Timeout:   15 * time.Second,
		Transport: transport,
	}

	return &client{
		client: httpClient,
		cfg: &Config{
			Host:    defaultHost,
			Format:  JsonFormat,
			Lang:    Ru,
			Version: apiVersion,
		},
		keys: keys,
	}
}

type SchedulesRequest struct {
	Station       string        //
	Time          time.Time     //
	TransportType TransportType //
	Offset        int
	Limit         int
}

type SchedulesResponse struct {
	Pagination        Pagination  `json:"pagination"`         // Информация о постраничном выводе найденных рейсов.
	Date              string      `json:"date"`               // Дата, на которую получен список рейсов.
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
	q.Set("apikey", c.key())
	q.Set("format", c.cfg.Format.String())
	q.Set("lang", c.cfg.Lang.String())
	q.Set("transport_type", req.TransportType.String())
	q.Set("station", req.Station)
	q.Set("date", req.Time.Format(dateFormat))

	if req.Offset != 0 {
		q.Set("offset", strconv.Itoa(req.Offset))
	}
	if req.Limit != 0 {
		q.Set("limit", strconv.Itoa(req.Limit))
	}

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
	q.Set("apikey", c.key())
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
	From   string
	To     string
	Date   time.Time
	Offset int
	Limit  int
}

type SearchResponse struct {
	Pagination       Pagination    `json:"pagination"`
	IntervalSegments []interface{} `json:"interval_segments"`
	Segments         []Segment     `json:"segments"`
	Search           interface{}   `json:"search"`
}

func (c *client) Search(ctx context.Context, req SearchRequest) (*SearchResponse, error) {
	if req.From == "" || req.To == "" {
		return nil, errors.New("one of required request param are missing")
	}

	u := url.URL{
		Scheme: scheme,
		Host:   c.cfg.Host,
		Path:   c.cfg.Version + "/search/",
	}

	q := u.Query()
	q.Set("apikey", c.key())
	q.Set("format", c.cfg.Format.String())
	q.Set("lang", c.cfg.Lang.String())
	q.Set("from", req.From)
	q.Set("to", req.To)

	if !req.Date.IsZero() {
		q.Set("date", req.Date.Format(dateFormat))
	}

	if req.Offset != 0 {
		q.Set("offset", strconv.Itoa(req.Offset))
	}
	if req.Limit != 0 {
		q.Set("limit", strconv.Itoa(req.Limit))
	}

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
	Days      string     `json:"days"`
	Stops     []Stop     `json:"stops"`
	Transport *Transport `json:"transport_subtype"`
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
	q.Set("apikey", c.key())
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

type NearestStationsRequest struct {
	Lat      float64
	Lng      float64
	Distance int
	Offset   int
	Limit    int
}

type NearestStationsResponse struct {
	Pagination Pagination `json:"pagination"`
	Stations   []Station  `json:"stations"`
}

func (c *client) NearestStations(ctx context.Context, req NearestStationsRequest) (*NearestStationsResponse, error) {
	if req.Lat == 0 || req.Lng == 0 || req.Distance == 0 {
		return nil, fmt.Errorf("unable to require params")
	}

	u := url.URL{
		Scheme: scheme,
		Host:   c.cfg.Host,
		Path:   c.cfg.Version + "/nearest_stations/",
	}

	q := u.Query()
	q.Set("apikey", c.key())
	q.Set("format", c.cfg.Format.String())
	q.Set("lang", c.cfg.Lang.String())
	q.Set("lat", strconv.FormatFloat(req.Lat, 'f', -1, 64))
	q.Set("lng", strconv.FormatFloat(req.Lng, 'f', -1, 64))
	q.Set("distance", strconv.Itoa(req.Distance))

	if req.Offset != 0 {
		q.Set("offset", strconv.Itoa(req.Offset))
	}
	if req.Limit != 0 {
		q.Set("limit", strconv.Itoa(req.Limit))
	}

	u.RawQuery = q.Encode()

	var resp NearestStationsResponse
	if err := c.get(ctx, u.String(), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

type NearestCityRequest struct {
	Lat      float64
	Lng      float64
	Distance int
	Offset   int
	Limit    int
}

type NearestCityResponse struct {
	Distance     float64 `json:"distance"`
	Code         string  `json:"code"`
	Title        string  `json:"title"`
	PopularTitle string  `json:"popular_title"`
	ShortTitle   string  `json:"short_title"`
	Lat          float64 `json:"lat"`
	Lng          float64 `json:"lng"`
	Type         string  `json:"type"`
}

func (c *client) NearestCity(ctx context.Context, req NearestCityRequest) (*NearestCityResponse, error) {
	if req.Lat == 0 || req.Lng == 0 {
		return nil, fmt.Errorf("unable to require params")
	}

	u := url.URL{
		Scheme: scheme,
		Host:   c.cfg.Host,
		Path:   c.cfg.Version + "/nearest_settlement/",
	}

	q := u.Query()
	q.Set("apikey", c.key())
	q.Set("format", c.cfg.Format.String())
	q.Set("lang", c.cfg.Lang.String())
	q.Set("lat", strconv.FormatFloat(req.Lat, 'f', -1, 64))
	q.Set("lng", strconv.FormatFloat(req.Lng, 'f', -1, 64))
	if req.Distance != 0 {
		q.Set("distance", strconv.Itoa(req.Distance))
	}

	if req.Offset != 0 {
		q.Set("offset", strconv.Itoa(req.Offset))
	}
	if req.Limit != 0 {
		q.Set("limit", strconv.Itoa(req.Limit))
	}

	u.RawQuery = q.Encode()

	var resp NearestCityResponse
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
		if httpResp.StatusCode == http.StatusTooManyRequests {
			err := c.nextKey()
			if err != nil {
				return err
			}
		}
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

func (c *client) key() string {
	return c.keys[c.idxKey]
}

func (c *client) nextKey() error {
	if c.idxKey == len(c.keys)-1 {
		return errors.New("pool api keys is empty")
	}

	c.idxKey = c.idxKey + 1

	return nil
}
