package model

import (
	"github.com/easy-say-api/src/config"
	"time"
)

type Message struct {
	Id           string    `json:"id"`
	Content      string    `json:"content"`
	LikeCount    int       `json:"like_count"`
	CommentCount int       `json:"comment_count"`
	UserId       string    `json:"-"`
	CreatedAt    time.Time `json:"created_at" xorm:"created"`
	UpdatedAt    time.Time `json:"updated_at" xorm:"updated"`
	IsEnable     bool      `json:"is_enable"`

	User User `json:"user" xorm:"extends"`

	Liked bool `json:"liked" sql:"-"`
}

type SimpleMessage struct {
	Id           string    `json:"id"`
	Content      string    `json:"content"`
	LikeCount    int       `json:"like_count"`
	CommentCount int       `json:"comment_count"`
	UserId       string    `json:"user_id"`
	CreatedAt    time.Time `json:"created_at" xorm:"created"`
	UpdatedAt    time.Time `json:"updated_at" xorm:"updated"`
}

func (u Message) TableName() string       { return "messages" }
func (u SimpleMessage) TableName() string { return "messages" }

func (m *Message) Get() error {
	_, err := config.Global.DB.Join("left", "users", "users.id=messages.user_id").Get(m)
	return err
}

func (m Message) GetList(limit int, offset int) ([]Message, error) {
	list := make([]Message, 0)
	err := config.Global.DB.
		Where("messages.is_enable=true").
		Join("left", "users", "users.id=messages.user_id").
		OrderBy("messages.created_at desc").
		Limit(limit, offset).
		Find(&list)

	return list, err
}

func (m *SimpleMessage) Add() error {
	_, err := config.Global.DB.
		InsertOne(m)
	return err
}
