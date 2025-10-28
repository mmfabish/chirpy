package auth

import (
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
