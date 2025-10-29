package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	token := headers.Get("Authorization")
	if token == "" {
		return "", errors.New("authorization header missing from request")
	}

	if !strings.HasPrefix(token, "Bearer ") {
		return "", errors.New("malformed Authorization header found in request")
	}

	return strings.TrimPrefix(token, "Bearer "), nil
}

func MakeRefreshToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(token), nil
}
