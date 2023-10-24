package repr

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
)

func hasherHex(data []byte) string {
	if true { // TODO: parse hash type from config
		digest := sha1.Sum(data)
		return hex.EncodeToString(digest[:])
	} else {
		digest := sha256.Sum256(data)
		return hex.EncodeToString(digest[:])
	}
}
