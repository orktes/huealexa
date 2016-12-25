package zway

import "time"

type Config struct {
	Hostname    string
	Port        string
	Username    string
	Password    string
	PollTimeout time.Duration
}
