package gitok_cat

import (
	"os"
	"path/filepath"

	"github.com/magnickolas/gitok/constants"
	"github.com/magnickolas/gitok/repr"
)

func CatObject(digest string) ([]byte, error) {
	dirName, fileName := digest[:2], digest[2:]
	obj, err := readObject(dirName, fileName)
	if err != nil {
		return nil, err
	}
	return obj.Representation(), nil
}

func readObject(dirName, fileName string) (repr.Object, error) {
	path := filepath.Join(constants.Git, constants.Objects, dirName, fileName)
	compressed, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return repr.ParseObject(compressed)
}
