package model

import (
	"errors"
	"time"
)

var (
	ErrChatNotFound = errors.New("chat not found")
)

type Chat struct {
	Usernames []string
}

type Message struct {
	From      string
	Text      string
	ChatID    int64
	Timestamp time.Time
}
