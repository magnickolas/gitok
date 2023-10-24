package gitok_hash

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/magnickolas/gitok/constants"
	"github.com/magnickolas/gitok/repr"
)

// Optionally saves the blob and return its key
func ProcessBlob(r io.Reader, save bool) (string, error) {
	blob, err := repr.NewBlob(r)
	if err != nil {
		return "", err
	}
	if save {
		err = saveObject(blob)
		if err != nil {
			return "", err
		}
	}
	return blob.Hash(), nil
}

func saveObject(obj repr.Object) error {
	b, err := obj.Compressed()
	if err != nil {
		return err
	}
	digest := obj.Hash()
	dirName, fileName := digest[:2], digest[2:]
	return writeObject(dirName, fileName, b)
}

func writeObject(dirName, fileName string, b []byte) error {
	objDirPath := filepath.Join(constants.Git, constants.Objects, dirName)
	_, err := os.Stat(objDirPath)
	if errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(objDirPath, os.ModePerm)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	objPath := filepath.Join(objDirPath, fileName)
	_, err = os.Stat(objPath)
	if errors.Is(err, os.ErrNotExist) {
		err = os.WriteFile(objPath, b, 0444)
		if err != nil {
			return err
		}
		return nil
	} else {
		return err
	}
}
