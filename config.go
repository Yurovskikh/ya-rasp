package yandex

import "time"

type Config struct {
	Host    string
	Format  format
	Lang    lang
	Version string
	Timeout time.Duration
}
