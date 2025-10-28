package auth

import (
	"testing"
)

func TestCheckPasswordHash(t *testing.T) {
	password := "Th1s!sAT3st!25"

	hash, err := HashPassword(password)
	if err != nil {
		t.Error(err)
	}

	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Error(err)
	}

	if !match {
		t.Errorf("Password and hash do not match: %s != %s", password, hash)
	}
}
