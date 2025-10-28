package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	userId, err := uuid.Parse("0940408c-4a2d-42ab-8680-0e2ed5bf899d")
	if err != nil {
		t.Error(err)
	}

	token, err := MakeJWT(userId, "TestMakeJWT", time.Millisecond*1000)
	if err != nil {
		t.Error(err)
	}

	if len(token) == 0 {
		t.Errorf("Invalid token: %s", token)
	}
}

func TestValidateJWT(t *testing.T) {
	userId, err := uuid.Parse("0940408c-4a2d-42ab-8680-0e2ed5bf899d")
	if err != nil {
		t.Error(err)
	}

	token, err := MakeJWT(userId, "TestMakeJWT", time.Millisecond*1000)
	if err != nil {
		t.Error(err)
	}

	subject, err := ValidateJWT(token, "TestMakeJWT")
	if err != nil {
		t.Error(err)
	}

	if subject != userId {
		t.Errorf("Subject does not match user ID: %s != %s", subject, userId)
	}
}
