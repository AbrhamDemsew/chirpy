package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name      string
		header    string
		wantToken string
		wantErr   bool
	}{
		{
			name:      "valid bearer token",
			header:    "Bearer token-string",
			wantToken: "token-string",
		},
		{
			name:      "extra whitespace",
			header:    "  Bearer    token-string  ",
			wantToken: "token-string",
		},
		{
			name:    "missing header",
			wantErr: true,
		},
		{
			name:    "wrong scheme",
			header:  "Basic token-string",
			wantErr: true,
		},
		{
			name:    "missing token",
			header:  "Bearer",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := http.Header{}
			if tt.header != "" {
				headers.Set("Authorization", tt.header)
			}

			token, err := GetBearerToken(headers)
			if tt.wantErr {
				if err == nil {
					t.Fatal("GetBearerToken returned no error")
				}
				return
			}

			if err != nil {
				t.Fatalf("GetBearerToken returned an error: %v", err)
			}

			if token != tt.wantToken {
				t.Errorf("got token %q, want %q", token, tt.wantToken)
			}
		})
	}
}

func TestMakeJWTAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	tokenString, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT returned an error: %v", err)
	}

	gotUserID, err := ValidateJWT(tokenString, secret)
	if err != nil {
		t.Fatalf("ValidateJWT returned an error: %v", err)
	}

	if gotUserID != userID {
		t.Errorf("got user ID %v, want %v", gotUserID, userID)
	}
}

func TestValidateJWTRejectsExpiredToken(t *testing.T) {
	tokenString, err := MakeJWT(uuid.New(), "test-secret", -time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT returned an error: %v", err)
	}

	if _, err := ValidateJWT(tokenString, "test-secret"); err == nil {
		t.Error("ValidateJWT accepted an expired token")
	}
}

func TestValidateJWTRejectsWrongSecret(t *testing.T) {
	tokenString, err := MakeJWT(uuid.New(), "correct-secret", time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT returned an error: %v", err)
	}

	if _, err := ValidateJWT(tokenString, "wrong-secret"); err == nil {
		t.Error("ValidateJWT accepted a token signed with a different secret")
	}
}
