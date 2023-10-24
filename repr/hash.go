package repr

import (
	"bytes"
	"compress/zlib"
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

func compress(raw []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	zw := zlib.NewWriter(buf)
	_, err := zw.Write(raw)
	if err != nil {
		return nil, err
	}
	// close so that the checksum is written to the buffer
	err = zw.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
