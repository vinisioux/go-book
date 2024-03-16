package security

import "golang.org/x/crypto/bcrypt"

func Hash(content string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(content), bcrypt.DefaultCost)
}

func CheckHash(hashContent, content string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashContent), []byte(content))
}
