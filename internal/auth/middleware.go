package auth

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/artemmarkaryan/exlex-backend/pkg/tokenizer"
	"github.com/cristalhq/jwt/v5"
)

const authHeader = "Authorization"
const keyToken = "auth_token"

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), keyToken, r.Header.Get(authHeader))
			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}

const keyClaim = "auth_claim"

func DirectiveAuthorized(ctx context.Context, _ interface{}, next graphql.Resolver) (res interface{}, err error) {
	rawToken, ok := ctx.Value(keyToken).(string)
	if !ok || rawToken == "" {
		return nil, ErrUnauthenticated
	}

	rawTokenParts := strings.Split(rawToken, " ")
	if len(rawTokenParts) != 2 || rawTokenParts[0] != "Bearer" {
		log.Printf("error: bad auth header: %q: %v", rawToken, err)
		return nil, ErrUnauthenticated
	}

	t, err := jwt.Parse([]byte(rawTokenParts[1]), tokenizer.FromContext(ctx).Verifier())
	if err != nil {
		log.Printf("error: parse token: %v", err)
		return nil, ErrUnauthenticated
	}

	var claim = Claim{}
	err = t.DecodeClaims(&claim)
	if err != nil {
		log.Printf("error: decode claims: %v", err)
		return nil, ErrUnauthenticated
	}

	return next(context.WithValue(ctx, keyClaim, claim))
}

func FromContext(ctx context.Context) (Claim, error) {
	c, ok := ctx.Value(keyClaim).(Claim)
	if !ok {
		return Claim{}, ErrUnauthenticated
	}

	return c, nil
}
