package utils

import (
	"crypto/md5"
	"fmt"
	"io"
)

func HashPassword(password string) string {
	passwordHash := md5.New()
	io.WriteString(passwordHash, password)
	return fmt.Sprintf("%x", passwordHash.Sum(nil))
}
