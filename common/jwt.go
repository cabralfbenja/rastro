package common

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/gommon/log"
	"os"
	"rastro/internal/models"
	"time"
)

type JWTClaims struct {
	ID uint `json:"id"`
	jwt.RegisteredClaims
}

func GenerateJWT(user models.UserModel) (*string, *string, error) {
	userClaims := JWTClaims{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	signedAccessToken, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, nil, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
		},
	})
	signedRefreshToken, err := refreshToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return nil, nil, err
	}

	return &signedAccessToken, &signedRefreshToken, nil
}

func ParseJWTSignedToken(signedAccessToken string) (*JWTClaims, error) {
	parsedJWTAccessToken, err := jwt.ParseWithClaims(signedAccessToken, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		log.Error(err)
		return nil, err
	} else if claims, ok := parsedJWTAccessToken.Claims.(*JWTClaims); ok && parsedJWTAccessToken.Valid {
		return claims, nil
	} else {
		return nil, errors.New("Unknown claims type, cannot proceed")
	}
}

func (claims *JWTClaims) IsExpired() bool {
	if claims.ExpiresAt == nil {
		return false
	}
	return claims.ExpiresAt.Time.Before(time.Now())
}
