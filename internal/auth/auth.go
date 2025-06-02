package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthInterface interface {
	GenerateToken(userID string, maxAccess int64, IntervalAccess int64, duration time.Duration) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
}

// Auth provides methods for generating and validating JWT tokens

type Auth struct {
	AuthInterface
	jwtSecret []byte
}

func NewAuth(secret []byte) *Auth {
	return &Auth{
		jwtSecret: secret,
	}
}

type Claims struct {
	MaxAccess      int64 `json:"max_access"`
	IntervalAccess int64 `json:"interval_access"`
	jwt.RegisteredClaims
}

// GenerateToken creates a JWT with a max access count and expiration
func (auth *Auth) GenerateToken(userID string, maxAccess int64, IntervalAccess int64, duration time.Duration) (string, error) {
	claims := Claims{
		MaxAccess:      maxAccess,
		IntervalAccess: IntervalAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(auth.jwtSecret)
}

// ValidateToken parses and validates the JWT, returning claims if valid
func (auth *Auth) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return auth.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
