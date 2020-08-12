package model

import (
	"github.com/easy-say-api/src/config"
	"time"
)

type User struct {
	Id         string `json:"id"`
	Nickname   string `json:"nickname"`
	Avatar     string `json:"avatar"`
	Username   string `json:"-"`
	Password   string `json:"-"`
	BriefIntro string `json:"brief_intro"`
	GithubId   int    `json:"-"`
	Email      string `json:"-"`

	CreatedAt time.Time `json:"-" xorm:"created"`
	UpdatedAt time.Time `json:"-" xorm:"updated"`
}

func (u User) TableName() string { return "users" }

func (u *User) Login() error {
	_, err := config.Global.DB.Get(u)
	return err
}

func (u *User) Get() error {
	_, err := config.Global.DB.Where("id=?", u.Id).Get(u)
	return err
}

func (u *User) GetByGithubId() error {
	_, err := config.Global.DB.Where("github_id=?", u.GithubId).Get(u)
	return err
}

func (u *User) Add() error {
	_, err := config.Global.DB.InsertOne(u)
	return err
}
