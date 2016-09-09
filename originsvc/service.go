package main

import "errors"

// OriginService provides operations on data from originating alert
type OriginService interface {
	ProcessAlert(string) (string, error)
}

type originService struct{}

func (originService) ProcessAlert(sid string) (string, error) {
	if sid == "" {
		return "", ErrEmpty
	}
	return sid, nil
}

// ErrEmpty is returned when an input string is empty.
var ErrEmpty = errors.New("empty string")
