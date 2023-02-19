package services

import "errors"

var (
	ErrLinkNotValid  = errors.New("link not valid")
	ErrLinkIsExist   = errors.New("link is exist")
	ErrLinkIsDeleted = errors.New("link is deleted")
)
