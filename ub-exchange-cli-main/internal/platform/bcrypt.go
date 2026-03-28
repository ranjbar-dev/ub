package platform

import "golang.org/x/crypto/bcrypt"

// PasswordEncoder provides bcrypt-based password hashing and verification.
type PasswordEncoder interface {
	// CompareHashAndPassword compares a bcrypt hashed password with a plaintext
	// candidate. Returns nil on success or an error if they do not match.
	CompareHashAndPassword(hashedPassword string, password string) error
	// GenerateFromPassword returns the bcrypt hash of the given plaintext password
	// using a cost factor of 12.
	GenerateFromPassword(password string) ([]byte, error)
}

type passwordEncoder struct {
}

func (pe *passwordEncoder) CompareHashAndPassword(hashedPassword string, password string) error {
	hashedPasswordBytes := []byte(hashedPassword)
	passwordBytes := []byte(password)
	return bcrypt.CompareHashAndPassword(hashedPasswordBytes, passwordBytes)
}

func (pe *passwordEncoder) GenerateFromPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 12)
}

func NewPasswordEncoder() PasswordEncoder {
	return &passwordEncoder{}
}
