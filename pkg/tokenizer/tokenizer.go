package tokenizer

import (
	"errors"
	"fmt"

	"github.com/cristalhq/jwt/v5"
)

// TODO: Make them expiraible

var (
	ErrNoAlgorithm = errors.New("no hashing algorithm provided in config")
	ErrNoSecret    = errors.New("no secret provided in config")
)

type Config struct {
	Algorithm jwt.Algorithm
	SecretKey string
}

type Tokenizer struct {
	builder  *jwt.Builder
	verifier *jwt.HSAlg
}

func MakeTokenizer(config Config) (f Tokenizer, err error) {
	if config.Algorithm == "" {
		return f, ErrNoAlgorithm
	}

	if config.SecretKey == "" {
		return f, ErrNoSecret
	}

	signer, err := jwt.NewSignerHS(config.Algorithm, []byte(config.SecretKey))
	f.builder = jwt.NewBuilder(signer)
	f.verifier, err = jwt.NewVerifierHS(config.Algorithm, []byte(config.SecretKey))
	if err != nil {
		err = fmt.Errorf("create tokenizer verifier: %w", err)
		return
	}

	return
}

func (t Tokenizer) NewToken(claims any) (*jwt.Token, error) {
	return t.builder.Build(claims)
}

func (t Tokenizer) VerifyToken(token *jwt.Token) (err error) {
	return t.verifier.Verify(token)
}

func (t Tokenizer) Parse(raw []byte, claims *any) error {
	return jwt.ParseClaims(raw, t.verifier, claims)
}
