package jwt

import (
	"time"

	"github.com/pascaldekloe/jwt"
)

type Tokenizer interface {
	Sign()
}

type TokenManager struct {
	config Config
}

func NewTokenizer(cfg Config) (t TokenManager, err error) {
	t.config, err = cfg, cfg.Validate()
	return
}

func (t TokenManager) Sign() {
	var c = new(jwt.Claims)
	c.Subject = ""
	c.Issued = jwt.NewNumericTime(time.Now())
	c.Expires = jwt.NewNumericTime(time.Now().Add(t.config.AccessExpires))
}
