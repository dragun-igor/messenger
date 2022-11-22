package model

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

type Message struct {
	Sender   string
	Receiver string
	Message  string
}

type AuthData struct {
	Login    string `validate:"min:4|max:20|regexp:[a-zA-z0-9]"`
	Password string `validate:"min:8"`
	Name     string `validate:"min:4|max:20|regexp:[a-zA-zа-яА-Я0-9]"`
}

type ValidationErrors []string

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

func Validate(x interface{}) (ValidationErrors, error) {
	ve := ValidationErrors{}
	v := reflect.ValueOf(x)
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		structField := t.Field(i)
		value := v.Field(i)
		tag := structField.Tag.Get("validate")
		if tag == "" {
			continue
		}
		tags := strings.Split(tag, "|")
		var err error
		switch value.Kind() {
		case reflect.String:
			ve, err = validateString(ve, structField.Name, tags, value.String())
		default:
			return nil, fmt.Errorf("unimplemented type: %s field: %s", value.Kind().String(), structField.Name)
		}
		if err != nil {
			return nil, err
		}
	}
	return ve, nil
}

func validateString(ve ValidationErrors, fieldName string, tags []string, value string) (ValidationErrors, error) {
	for _, tag := range tags {
		tagSplit := strings.Split(tag, ":")
		switch tagSplit[0] {
		case "min":
			tagValue, err := strconv.Atoi(tagSplit[1])
			if err != nil {
				return nil, err
			}
			ve = validateStringMinLen(ve, fieldName, value, tagValue)
		case "max":
			tagValue, err := strconv.Atoi(tagSplit[1])
			if err != nil {
				return nil, err
			}
			ve = validateStringMaxLen(ve, fieldName, value, tagValue)
		case "regexp":
			var err error
			ve, err = validateStringRegExp(ve, fieldName, value, tagSplit[1])
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unknown tag: %s field: %s", tagSplit[0], fieldName)
		}
	}
	return ve, nil
}

func validateStringMinLen(ve ValidationErrors, field, value string, tagValue int) ValidationErrors {
	counter := utf8.RuneCount([]byte(value))
	if counter < tagValue {
		ve = append(ve, fmt.Sprintf("%s must contain at least %d characters", field, tagValue))
	}
	return ve
}

func validateStringMaxLen(ve ValidationErrors, field, value string, tagValue int) ValidationErrors {
	counter := utf8.RuneCount([]byte(value))
	if counter > tagValue {
		ve = append(ve, fmt.Sprintf("%s must contain a maximum %d characters", field, tagValue))
	}
	return ve
}

func validateStringRegExp(ve ValidationErrors, field, value, tagValue string) (ValidationErrors, error) {
	ok, err := regexp.Match(tagValue, []byte(value))
	if err != nil {
		return nil, err
	}
	if !ok {
		ve = append(ve, fmt.Sprintf("%s does not match regular expression", field))
	}
	return ve, nil
}
