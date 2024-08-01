package web

import (
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestNewUserCrypto(t *testing.T) {
	password := "password"
	fromPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	err = bcrypt.CompareHashAndPassword(fromPassword, []byte(password))
	if err != nil {
		t.Errorf("password mismatch")
	}
}
