package search

import (
	"errors"
)

var ErrUnauthorized = errors.New("unauthorized")
var ErrNotFound = errors.New("not found")
