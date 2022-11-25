package service

import "github.com/dragun-igor/messenger/internal/pkg/model"

type authMatcher struct {
	model.AuthData
}

func NewAuthMatcher(data model.AuthData) authMatcher {
	return authMatcher{data}
}

func (a authMatcher) Matches(x interface{}) bool {
	a2, ok := x.(model.AuthData)
	if !ok {
		return false
	}
	return a2.Login == a.AuthData.Login && a2.Name == a.AuthData.Name && a2.IsPasswordCorrect(a.AuthData.Password)
}

func (a authMatcher) String() string {
	return ""
}
