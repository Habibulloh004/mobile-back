package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"mobilka/internal/models"

	"github.com/golang-jwt/jwt/v4"
)

// JWT secret key - auto-generated on application startup
// This is a placeholder and will be replaced with a generated key
var JWTSecret = []byte("auto_generated_secret_will_be_placed_here")

// InitJWTSecret generates a new JWT secret key
func InitJWTSecret() error {
	// Generate a random key for JWT token signing
	key := make([]byte, 32) // 256 bits
	_, err := rand.Read(key)
	if err != nil {
		return err
	}
	JWTSecret = key
	return nil
}

// Claims represents the JWT claims
type Claims struct {
	ID   int    `json:"id"`
	Role string `json:"role"` // "admin" or "superadmin"
	jwt.RegisteredClaims
}

// GenerateToken generates a new JWT token for the provided admin
func GenerateAdminToken(admin *models.Admin) (string, error) {
	// Create claims with admin info
	claims := Claims{
		ID:   admin.ID,
		Role: "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate the JWT
	tokenString, err := token.SignedString(JWTSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateSuperAdminToken generates a new JWT token for the provided super admin
func GenerateSuperAdminToken(superAdmin *models.SuperAdmin) (string, error) {
	// Create claims with super admin info
	claims := Claims{
		ID:   superAdmin.ID,
		Role: "superadmin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate the JWT
	tokenString, err := token.SignedString(JWTSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken parses and validates a JWT token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// GenerateSystemToken generates a random token for system authentication
func GenerateSystemToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(token), nil
}

// GenerateSmsToken generates a random token for SMS authentication
func GenerateSmsToken() (string, error) {
	token := make([]byte, 16)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(token), nil
}
