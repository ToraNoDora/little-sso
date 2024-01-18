package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/ToraNoDora/little-sso/sso/internal/src/domain/models"
)

// NewToken creates new JWT token for given user and app
func NewToken(user models.User, app models.App, duration time.Duration, hashingPermissions string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["app_id"] = app.ID
	claims["hash"] = hashingPermissions
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(token string, appSecret string) error {
	tokenParsed, err := jwt.Parse(
		token,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid signing method")
			}
			return []byte(appSecret), nil
		},
	)
	if err != nil {
		return fmt.Errorf("failed to parse token: %v", err.Error())
	}

	if tokenParsed.Valid {
		if claims, ok := tokenParsed.Claims.(jwt.MapClaims); ok {
			_ = claims["hash"].(string)
			expirationFloat := claims["exp"].(float64)
			if !ok {
				return fmt.Errorf("expiration not found in claims: %v", err.Error())
			}

			// Compare the expiration time with the current time
			expiration := time.Unix(int64(expirationFloat), 0)
			if time.Now().After(expiration) {
				return fmt.Errorf("token has expired: %v", err.Error())
			}

		} else {
			return fmt.Errorf("invalid claims: %v", err.Error())
		}
	} else {
		return fmt.Errorf("token is invalid: %v", err.Error())
	}

	return nil
}
