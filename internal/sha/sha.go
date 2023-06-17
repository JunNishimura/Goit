package sha

import (
	"encoding/hex"
	"errors"
	"regexp"
)

type SHA1 []byte

var (
	sha1Regexp     = regexp.MustCompile("[0-9a-f]{40}")
	ErrInvalidHash = errors.New("invalid hash")
)

func (s SHA1) String() string {
	return hex.EncodeToString(s)
}

func (s SHA1) Compare(other SHA1) bool {
	return s.String() == other.String()
}

func ReadHash(hashString string) (SHA1, error) {
	if ok := sha1Regexp.MatchString(hashString); !ok {
		return nil, ErrInvalidHash
	}
	hash, err := hex.DecodeString(hashString)
	if err != nil {
		return nil, ErrInvalidHash
	}
	return hash, nil
}
