package model

import "golang.org/x/crypto/bcrypt"

type Message struct {
	Sender   string
	Receiver string
	Message  string
}

type User struct {
	Login    string
	Password string
	Name     string
}

func (u *User) SetHashByPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

func (u *User) IsPasswordCorrect(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}
