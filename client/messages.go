package main

import (
	"time"

	"github.com/dragun-igor/messenger/messengerpb"
)

type Message struct {
	Id       int64             `json:"id"`
	Time     time.Time         `json:"time"`
	Sender   *messengerpb.User `json:"sender"`
	Receiver *messengerpb.User `json:"receiver"`
	Msg      string            `json:"msg"`
}

type User struct {
	Id         int64      `json:"id"`
	FirstName  string     `json:"first_name"`
	SecondName string     `json:"second_name"`
	Messages   []*Message `json:"messages"`
}

type List map[int64]*User

// func GetList() {
// 	list := make(List)
// 	GetAl
// }
