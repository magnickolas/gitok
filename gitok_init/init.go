package gitok_init

import (
	"fmt"
	"os"
	"path/filepath"
)

const git = ".git"
const head = "HEAD"
const objects = "objects"
const refs = "refs"

const refFormat = "ref: refs/heads/%v"

func InitRepo(initBranch string) error {
	err := os.Mkdir(git, os.ModePerm)
	if err != nil {
		return err
	}
	err = os.Mkdir(filepath.Join(git, objects), os.ModePerm)
	if err != nil {
		return err
	}
	err = os.Mkdir(filepath.Join(git, refs), os.ModePerm)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(git, head),
		[]byte(fmt.Sprintf(refFormat, initBranch)), 0644)
	if err != nil {
		return err
	}
	return nil
}
