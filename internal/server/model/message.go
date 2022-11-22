package model

import "golang.org/x/crypto/bcrypt"

type Message struct {
	Sender   string
	Receiver string
	Message  string
}

type AuthData struct {
	Login    string
	Password string
	Name     string
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
