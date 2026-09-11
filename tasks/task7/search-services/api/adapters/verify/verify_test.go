package verify

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
)

func TestTokenVerification(t *testing.T) {
	for _, tc := range []struct {
		name            string
		ttl             time.Duration
		verificationKey string
		wantError       bool
	}{
		{"valid", time.Minute, "test-secret", false},
		{"expired", -time.Minute, "test-secret", true},
		{"wrong key", time.Minute, "wrong-secret", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			token, err := CreateToken("superuser", "test-secret", tc.ttl)
			if err != nil {
				t.Fatal(err)
			}
			if err := VerifyToken(token, tc.verificationKey); (err != nil) != tc.wantError {
				t.Fatalf("VerifyToken() error = %v, wantError = %v", err, tc.wantError)
			}
		})
	}
}

func TestRejectOtherSigningAlgorithm(t *testing.T) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS384, jwt.MapClaims{
		"sub": "superuser",
		"exp": time.Now().Add(time.Minute).Unix(),
	}).SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	if err := VerifyToken(token, "test-secret"); err == nil {
		t.Fatal("accepted token signed with HS384")
	}
}
