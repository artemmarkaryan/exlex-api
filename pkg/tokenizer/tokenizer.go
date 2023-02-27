package tokenizer

import (
	"context"
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
	bld *jwt.Builder
	ver *jwt.HSAlg
}

func MakeTokenizer(config Config) (f Tokenizer, err error) {
	if config.Algorithm == "" {
		return f, ErrNoAlgorithm
	}
	if config.SecretKey == "" {
		return f, ErrNoSecret
	}

	signer, err := jwt.NewSignerHS(config.Algorithm, []byte(config.SecretKey))
	f.bld = jwt.NewBuilder(signer)
	f.ver, err = jwt.NewVerifierHS(config.Algorithm, []byte(config.SecretKey))
	if err != nil {
		err = fmt.Errorf("create tokenizer verifier: %w", err)
		return
	}

	return
}

func (t Tokenizer) NewToken(claims any) (*jwt.Token, error)  { return t.bld.Build(claims) }
func (t Tokenizer) Verifier() jwt.Verifier                   { return t.ver }
func (t Tokenizer) VerifyToken(token *jwt.Token) (err error) { return t.ver.Verify(token) }
func (t Tokenizer) Parse(raw []byte, claims *any) error      { return jwt.ParseClaims(raw, t.ver, claims) }

const tokenizerKey = "tokenizer"

func Propagate(ctx context.Context, t Tokenizer) context.Context {
	return context.WithValue(ctx, tokenizerKey, t)
}

func FromContext(ctx context.Context) Tokenizer {
	t, ok := ctx.Value(tokenizerKey).(Tokenizer)
	if !ok {
		panic("no tokenizer in " + tokenizerKey)
	}
	return t
}
