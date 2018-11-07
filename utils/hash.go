package utils

import (
	"crypto/md5"
	"fmt"
	"io"
)

func GeneratePasswordHash(password string) string {
	passwordHash := md5.New()
	io.WriteString(passwordHash, password)
	return fmt.Sprintf("%x", passwordHash.Sum(nil))
}

func CompareHashAndPassword(hash, password string) bool {
	return hash == GeneratePasswordHash(password)
}
