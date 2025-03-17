package hashing

import (
	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct{}

func New() *BcryptHasher {
	return &BcryptHasher{}
}

func (b *BcryptHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (b *BcryptHasher) Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
