package model

import (
	"golang.org/x/crypto/bcrypt"
)

type Message struct {
	Sender   string
	Receiver string
	Message  string
}

type AuthData struct {
	Login    string `validate:"min:1|max:32|regexp:[a-zA-Z0-9]"`
	Password string `validate:"min:8|max:32|regexp:[a-zA-Z0-9]"`
	Name     string `validate:"min:1|max:32|regexp:[a-zA-zа-яА-Я]"`
}

func (a *AuthData) SetHashByPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	a.Password = string(hash)
	return nil
}

func (a *AuthData) IsPasswordCorrect(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password)) == nil
}
