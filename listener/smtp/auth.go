package smtp

import (
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

func NewAuth(enabled bool, users []*AuthUser) *Auth {
	return &Auth{
		Enabled: enabled,
		Users:   users,
	}
}

func NewAuthUser(email, passwordHash string) *AuthUser {
	return &AuthUser{Email: email, PasswordHash: passwordHash}
}

func (au *AuthUser) Authenticate(username, password string) bool {
	return au.Email == username && bcrypt.CompareHashAndPassword([]byte(au.PasswordHash), []byte(password)) == nil
}
