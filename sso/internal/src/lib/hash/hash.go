package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/ToraNoDora/little-sso/sso/internal/src/domain/models"
)

func HashPassword(password string) ([]byte, error) {
	// Generate a salt with a cost factor of 10
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

func VerifyPassword(password string, hashedPassword []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
}

func HashingPermissions(prms []models.Permission) (string, error) {
	hasher := sha256.New()
	objStr := fmt.Sprintf("%+v", prms)
	data := []byte(objStr)

	_, err := hasher.Write(data)
	if err != nil {
		return "", err
	}

	hash := hasher.Sum(nil)
	hashString := hex.EncodeToString(hash)

	return hashString, nil
}
