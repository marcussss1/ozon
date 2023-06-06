package hash

import (
	"encoding/base64"
	"golang.org/x/crypto/argon2"
	"regexp"
)

var salt = []byte{0x5C, 0x72, 0x69, 0x67, 0x61, 0x64, 0x65}

func Hash(originalUrl string) string {
	abbreviatedUrlHash := argon2.IDKey([]byte(originalUrl), salt, 1, 1, 1, 32)
	abbreviatedUrl := base64.StdEncoding.EncodeToString(abbreviatedUrlHash)

	re := regexp.MustCompile("[^a-zA-Z0-9_]")
	abbreviatedUrl = re.ReplaceAllString(abbreviatedUrl, "")

	return abbreviatedUrl[:10]
}
