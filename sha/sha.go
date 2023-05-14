package sha

import (
	"encoding/hex"
	"fmt"
	"regexp"
)

type SHA1 []byte

var (
	sha1Regexp = regexp.MustCompile("[0-9a-f]{40}")
)

func (sha1 SHA1) String() string {
	return hex.EncodeToString(sha1)
}

func ReadHash(hashString string) (SHA1, error) {
	if ok := sha1Regexp.MatchString(hashString); !ok {
		return nil, fmt.Errorf("invalid hash")
	}
	hash, err := hex.DecodeString(hashString)
	if err != nil {
		return nil, fmt.Errorf("fail to decode hash string: %v", err)
	}
	return hash, nil
}
