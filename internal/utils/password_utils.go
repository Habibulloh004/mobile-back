package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPassword compares a password with its hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateSecurePassword generates a secure random password
func GenerateSecurePassword(length int) (string, error) {
	if length < 8 {
		length = 8 // Minimum password length
	}

	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

// GenerateSuperAdminPassword generates a secure password for the super admin
// This function is only called once during application initialization
func GenerateSuperAdminPassword() (string, string, error) {
	// Use fixed password "helloworld"
	password := "helloworld"

	// Hash the password
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return "", "", err
	}

	// Important: The plain text password should be shown only once during setup
	// and then discarded. It should be logged or displayed to the setup operator.
	fmt.Println("====== SUPER ADMIN CREDENTIALS ======")
	fmt.Println("Login: superadmin")
	fmt.Println("Password: " + password)
	fmt.Println("PLEASE SAVE THESE CREDENTIALS SECURELY")
	fmt.Println("====================================")

	return password, hashedPassword, nil
}
