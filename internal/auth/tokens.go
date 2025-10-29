package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func parseAuthorizationHeader(headers http.Header, prefix string) (string, error) {
	authString := headers.Get("Authorization")
	if authString == "" {
		return "", errors.New("authorization header missing from request")
	}

	authString = strings.TrimSpace(authString)

	if !strings.HasPrefix(authString, prefix) {
		return "", fmt.Errorf("malformed Authorization header: missing %s prefix", prefix)
	}

	return strings.TrimPrefix(authString, prefix), nil
}

func GetApiKey(headers http.Header) (string, error) {
	return parseAuthorizationHeader(headers, "ApiKey ")
}

func GetBearerToken(headers http.Header) (string, error) {
	return parseAuthorizationHeader(headers, "Bearer ")
}

func MakeRefreshToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(token), nil
}
