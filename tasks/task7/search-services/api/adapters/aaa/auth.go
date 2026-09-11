package aaa

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"yadro.com/course/api/adapters/verify"
	"yadro.com/course/api/core"
)

const secretKey = "something secret here" // token sign key
const adminRole = "superuser"             // token subject

// Authentication, Authorization, Accounting
type AAA struct {
	users    map[string]string
	tokenTTL time.Duration
	log      *slog.Logger
}

func New(tokenTTL time.Duration, log *slog.Logger) (AAA, error) {
	const adminUser = "ADMIN_USER"
	const adminPass = "ADMIN_PASSWORD"
	user, ok := os.LookupEnv(adminUser)
	if !ok {
		return AAA{}, fmt.Errorf("could not get admin user from enviroment")
	}
	password, ok := os.LookupEnv(adminPass)
	if !ok {
		return AAA{}, fmt.Errorf("could not get admin password from enviroment")
	}

	return AAA{
		users:    map[string]string{user: password},
		tokenTTL: tokenTTL,
		log:      log,
	}, nil
}

func (a AAA) Login(name, password string) (string, error) {
	if password1, ok := a.users[name]; ok && password1 == password {
		token, err := verify.CreateToken(adminRole, secretKey, a.tokenTTL)

		if err != nil {
			a.log.Error("aaa login", "error", err)
			return "", err
		}
		return token, nil
	}
	return "", core.ErrNotFound
}

func (a AAA) Verify(tokenString string) error {
	err := verify.VerifyToken(tokenString, secretKey)

	if err != nil {
		a.log.Error("aaa verify", "error", err)
	}

	return err
}
