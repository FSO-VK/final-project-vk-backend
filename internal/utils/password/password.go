package password

import (
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	Encrypt(password string) (string, error)
	Compare(password, hashedPassword string) bool
}

type PasswordHasherProvider struct{}

func NewPasswordHasherProvider() *PasswordHasherProvider {
	return &PasswordHasherProvider{}
}

func (p *PasswordHasherProvider) Encrypt(password string) (string, error) {
	byteHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(byteHash), nil
}

func (p *PasswordHasherProvider) Compare(password, hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
