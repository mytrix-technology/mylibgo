package helper

import (
	"context"
	"fmt"
	"strings"
	"github.com/mytrix-technology/mylibgo/security/cryptography"
)

type AuthObject struct {
	Token     string
	Datetime  string
	Signature string
	UserID    int64
}

func (a *AuthObject) ToHeadersMap() map[string]string {
	return map[string]string{
		"Authorization": fmt.Sprintf("token=%s", a.Token),
		"signature":     a.Signature,
		"datetime":      a.Datetime,
	}
}

func (a *AuthObject) GenerateSignature(signatureKey string) {
	a.Signature = cryptography.HmacSHA256(a.Token+a.Datetime, signatureKey)
}

func AuthObjectFromContext(ctx context.Context) (*AuthObject, error) {
	auth, ok := ctx.Value(ContextKeyRequestAuthorization).(string)
	if !ok {
		return nil, fmt.Errorf("cannot find Authorization context")
	}

	datetime, ok := ctx.Value(ContextKeyRequestDatetime).(string)
	if !ok {
		return nil, fmt.Errorf("cannot find datetime auth context")
	}

	signature, ok := ctx.Value(ContextKeyRequestSignature).(string)
	if !ok {
		return nil, fmt.Errorf("cannot find signature auth context")
	}

	if auth == "" {
		return nil, fmt.Errorf("invalid Authorization context")
	}

	tokens := strings.Split(auth, "=")
	if len(tokens) < 2 {
		return nil, fmt.Errorf("invalid Authorization context")
	}

	if strings.TrimSpace(tokens[0]) != "token" || strings.TrimSpace(tokens[1]) == "" {
		return nil, fmt.Errorf("invalid Authorization context")
	}

	token := strings.TrimSpace(tokens[1])

	if datetime == "" {
		return nil, fmt.Errorf("invalid datetime context")
	}

	if signature == "" {
		return nil, fmt.Errorf("invalid signature context")
	}

	return &AuthObject{
		Token:     token,
		Datetime:  datetime,
		Signature: signature,
	}, nil
}
