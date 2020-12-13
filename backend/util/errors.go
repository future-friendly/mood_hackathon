package util

import "fmt"

type AlreadyExists struct {
	Model string
	Key string
}

func (e AlreadyExists) Error() string {
	return fmt.Sprintf("%s with such %s already exists", e.Model, e.Key)
}

type WrongCredentials struct {
	Login string
}

func (e WrongCredentials) Error() string {
	return fmt.Sprint("wrong credentials for %s", e.Login)
}