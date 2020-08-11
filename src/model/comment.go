package model

import (
	"github.com/easy-say-api/src/config"
	"time"
)

type Comment struct {
	Id        string    `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at" xorm:"created"`
	UpdatedAt time.Time `json:"updated_at" xorm:"updated"`
	IsEnable  bool      `json:"is_enable"`

	UserId    string `json:"-"`
	MessageId string `json:"-"`

	User User `json:"user" xorm:"extends"`
}

type SimpleComment struct {
	Id        string    `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at" xorm:"created"`
	IsEnable  bool      `json:"is_enable"`

	UserId    string `json:"-"`
	MessageId string `json:"-"`
}

func (c Comment) TableName() string       { return "comments" }
func (c SimpleComment) TableName() string { return "comments" }

func (c *Comment) Get() error {
	_, err := config.Global.DB.Where("id=?", c.Id).Get(c)
	return err
}

func (c Comment) GetListByMessage(messageId string) ([]Comment, error) {
	list := make([]Comment, 0)
	err := config.Global.DB.
		Where("comments.is_enable=true and message_id=?", messageId).
		Join("left", "users", "users.id=comments.user_id").
		OrderBy("comments.created_at desc").
		Find(&list)

	return list, err
}

func (c *SimpleComment) Add() error {
	session := config.Global.DB.NewSession()
	_, err := session.InsertOne(c)
	if err == nil {
		msg := Message{Id: c.MessageId}
		_, err = session.Cols("comment_count").Incr("comment_count").Where("id=?", c.MessageId).Update(msg)
	}
	if err != nil {
		session.Rollback()
	} else {
		err = session.Commit()
	}
	return err
}

func (c *Comment) Del() error {
	session := config.Global.DB.NewSession()
	_, err := session.Cols("is_enable").Where("id=? and user_id=?", c.Id, c.UserId).Update(c)
	if err == nil {
		msg := Message{Id: c.MessageId}
		_, err = session.Cols("comment_count").Decr("comment_count").Where("id=?", c.MessageId).Update(msg)
	}
	if err != nil {
		session.Rollback()
	} else {
		err = session.Commit()
	}
	return err
}
