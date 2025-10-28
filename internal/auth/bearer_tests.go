package auth

import (
	"fmt"
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	expectedTokenString := "M5CO2zprQTBwpFfXEaEX4ZIO55SS659TZvH3Uoq2M8NNvWSWabAmRZt9s+Yrcs08"

	headers := http.Header{}
	headers.Add("Authorization", fmt.Sprintf("Bearer %s", expectedTokenString))

	actualTokenString, err := GetBearerToken(headers)
	if err != nil {
		t.Error(err)
	}

	if actualTokenString != expectedTokenString {
		t.Errorf("Token values do not match: %s != %s", actualTokenString, expectedTokenString)
	}
}
