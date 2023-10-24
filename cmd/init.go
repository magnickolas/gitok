package cmd

import (
	"github.com/magnickolas/gitok/gitok_init"
	"github.com/spf13/cobra"
)

var (
	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize a repo",
		Run: func(cmd *cobra.Command, args []string) {
			if err := gitok_init.InitRepo(branchName); err != nil {
				panic(err)
			}
		},
	}
	branchName string
)

func init() {
	const defaultBranchName = "master"

	initCmd.Flags().
		StringVarP(&branchName, "initial-branch", "b", defaultBranchName, "initial branch name")
}
