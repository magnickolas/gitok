package fs

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/magnickolas/gitok/constants"
	"github.com/magnickolas/gitok/repr"
)

func getObjectDirPath(digest string) string {
	dirName := digest[:2]
	return filepath.Join(constants.Git, constants.Objects, dirName)
}

func getObjectFilePath(digest string) string {
	dirPath, fileName := getObjectDirPath(digest), digest[2:]
	return filepath.Join(dirPath, fileName)
}

func ReadObject(digest string) (repr.Object, error) {
	path := getObjectFilePath(digest)
	compressed, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return repr.ParseObject(compressed)
}

func WriteObject(o repr.Object) error {
	compressed, err := o.Compressed()
	if err != nil {
		return err
	}
	objDirPath := getObjectDirPath(o.Digest())
	_, err = os.Stat(objDirPath)
	if errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(objDirPath, os.ModePerm)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	objPath := getObjectFilePath(o.Digest())
	_, err = os.Stat(objPath)
	if errors.Is(err, os.ErrNotExist) {
		err = os.WriteFile(objPath, compressed, 0444)
		if err != nil {
			return err
		}
		return nil
	} else {
		return err
	}
}
