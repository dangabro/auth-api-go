package process

import (
	"crypto/md5"
	"fmt"
	"strings"
)

func HashPassword(pass string) string {
	trimmedSpaceProvided := strings.TrimSpace(pass)
	providedBytes := []byte(trimmedSpaceProvided)
	hashedProvidedPassword := md5.Sum(providedBytes)
	hashedHexa := fmt.Sprintf("%x", hashedProvidedPassword)

	return hashedHexa
}
