package gitok_init

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/magnickolas/gitok/constants"
)

func InitRepo(initBranch string) error {
	err := os.Mkdir(constants.Git, os.ModePerm)
	if err != nil {
		return err
	}
	err = os.Mkdir(filepath.Join(constants.Git, constants.Objects), os.ModePerm)
	if err != nil {
		return err
	}
	err = os.Mkdir(filepath.Join(constants.Git, constants.Refs), os.ModePerm)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(constants.Git, constants.Head),
		[]byte(fmt.Sprintf(constants.RefFormat, initBranch)), 0644)
	if err != nil {
		return err
	}
	return nil
}
