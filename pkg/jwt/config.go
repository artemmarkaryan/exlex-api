package jwt

import (
	"errors"
	"time"
)

type Config struct {
	AccessExpires  time.Duration
	RefreshExpires time.Duration
}

func (c Config) Validate() error {
	if c.AccessExpires > 0 &&
		c.RefreshExpires > 0 {
		return nil
	}

	return errors.New("invalid tokenizer config")
}
