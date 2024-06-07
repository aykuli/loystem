package user

import (
	"crypto/sha512"
	"encoding/hex"
)

type User struct {
	ID             int    `json:"id"`
	Login          string `json:"login"`
	HashedPassword string `json:"-"`
	Salt           string
}

func (u *User) SetHashedPassword(password string, salt []byte) {
	u.HashedPassword = u.buildHashedPassword(password, salt)
}
func (u *User) buildHashedPassword(password string, salt []byte) string {
	passwordBytes := []byte(password)
	sha512Hasher := sha512.New()

	passwordBytes = append(passwordBytes, salt...)
	// Write password bytes to the hasher
	sha512Hasher.Write(passwordBytes)
	// Get the SHA-512 hashed password
	hashedPasswordBytes := sha512Hasher.Sum(nil)
	// Convert the password to a hex string
	hashedPasswordHex := hex.EncodeToString(hashedPasswordBytes)

	return hashedPasswordHex
}

func (u *User) ValidatePassword(password string) bool {
	salt, _ := hex.DecodeString(u.Salt)
	givenPasswordHash := u.buildHashedPassword(password, salt)

	return givenPasswordHash == u.HashedPassword
}
