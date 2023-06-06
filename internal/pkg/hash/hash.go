package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
)

func Hash(originalUrl string) string {
	h := sha256.New()
	h.Write([]byte(originalUrl))

	hash := hex.EncodeToString(h.Sum(nil))
	re := regexp.MustCompile("[^a-zA-Z0-9_]")
	hash = re.ReplaceAllString(hash, "")

	return hash[:10]
}
