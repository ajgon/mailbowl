package smtp

import (
	"github.com/ajgon/mailbowl/config"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	Enabled bool
	Users   []*AuthUser
}

type AuthUser struct {
	Email        string
	PasswordHash string
}

func NewAuth(conf config.SMTPAuth) *Auth {
	users := make([]*AuthUser, 0)

	for _, user := range conf.Users {
		users = append(users, &AuthUser{Email: user.Email, PasswordHash: user.PasswordHash})
	}

	return &Auth{
		Enabled: conf.Enabled,
		Users:   users,
	}
}

func NewAuthUser(email, passwordHash string) *AuthUser {
	return &AuthUser{Email: email, PasswordHash: passwordHash}
}

func (au *AuthUser) Authenticate(username, password string) bool {
	return au.Email == username && bcrypt.CompareHashAndPassword([]byte(au.PasswordHash), []byte(password)) == nil
}
