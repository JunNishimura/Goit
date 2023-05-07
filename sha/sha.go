package sha

import "encoding/hex"

type SHA1 []byte

func (sha1 SHA1) String() string {
	return hex.EncodeToString(sha1)
}
