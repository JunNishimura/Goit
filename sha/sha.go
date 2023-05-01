package sha

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
)

func HashGenerator(content string) string {
	sha1 := sha1.New()
	io.WriteString(sha1, content)
	return hex.EncodeToString(sha1.Sum(nil))
}
