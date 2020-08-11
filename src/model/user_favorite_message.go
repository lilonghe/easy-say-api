package model

import (
	"github.com/easy-say-api/src/config"
	"time"
)

type UserFavoriteMessage struct {
	Id        int       `json:"id"`
	UserId    string    `json:"user_id"`
	MessageId string    `json:"message_id"`
	CreatedAt time.Time `json:"created_at" xorm:"created"`
	UpdatedAt time.Time `json:"updated_at" xorm:"updated"`
	IsEnable  bool      `json:"is_enable"`
}

func (u UserFavoriteMessage) TableName() string { return "user_favorite_message" }

func (u UserFavoriteMessage) CheckUserInList(ids []string, user_id string) (arr []string, err error) {
	objs := make([]UserFavoriteMessage, 0)
	err = config.Global.DB.Where("is_enable=true and user_id=?", user_id).In("message_id", ids).Find(&objs)
	if err == nil {
		for _, v := range objs {
			arr = append(arr, v.MessageId)
		}
	}
	return arr, err
}

func (u *UserFavoriteMessage) Add() error {
	session := config.Global.DB.NewSession()
	_, err := session.InsertOne(u)
	if err == nil {
		msg := Message{Id: u.MessageId}
		_, err = session.Cols("like_count").Incr("like_count").Where("id=?", u.MessageId).Update(msg)
	}
	if err != nil {
		err = session.Rollback()
	} else {
		err = session.Commit()
	}
	return err
}

func (u *UserFavoriteMessage) Del() error {
	session := config.Global.DB.NewSession()
	_, err := session.Cols("is_enable").Where("message_id=? and user_id=?", u.MessageId, u.UserId).Update(u)
	if err == nil {
		msg := Message{Id: u.MessageId}
		_, err = session.Cols("like_count").Decr("like_count").Where("id=?", u.MessageId).Update(msg)
	}
	if err != nil {
		err = session.Rollback()
	} else {
		err = session.Commit()
	}
	return err
}

func (u *UserFavoriteMessage) GetByUserAndMessage() error {
	_, err := config.Global.DB.Where("is_enable=true").Get(u)
	return err
}
