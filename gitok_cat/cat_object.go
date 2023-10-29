package gitok_cat

import (
	"github.com/magnickolas/gitok/fs"
	"github.com/magnickolas/gitok/repr"
)

func GetObjectType(digest string) (string, error) {
	o, err := fs.ReadObject(digest)
	if err != nil {
		return "", err
	}
	return o.Type(), nil
}

func CatObject(digest string) (string, error) {
	o, err := fs.ReadObject(digest)
	if err != nil {
		return "", err
	}
	return string(repr.StripObjectHeader(o)), nil
}

func PrettyCatObject(digest string) (string, error) {
	o, err := fs.ReadObject(digest)
	if err != nil {
		return "", err
	}
	return o.String(), nil
}
