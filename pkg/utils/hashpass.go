package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPass(pass []byte) string {
  hashed, err := bcrypt.GenerateFromPassword(pass, bcrypt.MinCost)
  if err != nil {
    log.Printf("Failed to generate password: %v", err)
    return ""
  }

  return string(hashed)
}
