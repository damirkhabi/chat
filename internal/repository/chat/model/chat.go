package model

import (
	"time"
)

type Chat struct {
	ID        int64    `db:"id"`
	Usernames []string `db:"usernames"`
}

type Message struct {
	ID        int64     `db:"id"`
	From      string    `db:"from"`
	Text      string    `db:"text"`
	ChatID    int64     `db:"chat_id"`
	Timestamp time.Time `db:"timestamp"`
}
