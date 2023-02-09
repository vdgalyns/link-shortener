package services

import "errors"

var (
	ErrURLNotValid = errors.New("url not valid")
	ErrURLIsExist  = errors.New("url is exist")
)
