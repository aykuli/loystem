package user

import (
	"crypto/sha512"
	"encoding/hex"
)

type User struct {
	ID             int
	Login          string
	HashedPassword string
}

func (u *User) SetHashedPassword(password string, salt string) {
	u.HashedPassword = u.buildHashedPassword(password, salt)
}
func (u *User) buildHashedPassword(password string, salt string) string {
	passwordBytes := []byte(password)
	sha512Hasher := sha512.New()
	passwordBytes = append(passwordBytes, []byte(salt)...)
	// Write password bytes to the hasher
	sha512Hasher.Write(passwordBytes)
	// Get the SHA-512 hashed password
	hashedPasswordBytes := sha512Hasher.Sum(nil)
	// Convert the password to a hex string
	hashedPasswordHex := hex.EncodeToString(hashedPasswordBytes)

	return hashedPasswordHex
}

func (u *User) ValidatePassword(password string, salt string) bool {
	givenPasswordHash := u.buildHashedPassword(password, salt)

	return givenPasswordHash == u.HashedPassword
}
