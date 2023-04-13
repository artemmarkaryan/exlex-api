package search

import (
	"errors"
)

var ErrUnauthorized = errors.New("unauthorized")
var ErrNotFound = errors.New("not found")
var ErrApplicationAlreadyExists = errors.New("application for this search already exists")
