package gitok_hash

import (
	"io"

	"github.com/magnickolas/gitok/fs"
	"github.com/magnickolas/gitok/repr"
)

// Optionally saves the blob and return its key
func ProcessBlob(r io.Reader, save bool) (string, error) {
	blob, err := repr.NewBlob(r)
	if err != nil {
		return "", err
	}
	if save {
		err = fs.WriteObject(blob)
		if err != nil {
			return "", err
		}
	}
	return blob.Digest(), nil
}
