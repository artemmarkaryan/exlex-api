package graph

import (
	"errors"
	"fmt"
)

var ErrBadUUID = errors.New("bad UUID")
var ErrNotFound = func(what string) error { return fmt.Errorf("%s not found", what) }
