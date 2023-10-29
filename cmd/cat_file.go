package cmd

import (
	"fmt"

	"github.com/magnickolas/gitok/gitok_cat"
	"github.com/spf13/cobra"
)

var (
	catFileCmd = &cobra.Command{
		Use:   "cat-file",
		Short: "Print object info",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var str string
			var err error
			if prettyPrint {
				str, err = gitok_cat.PrettyCatObject(args[0])
			} else {
				str, err = gitok_cat.CatObject(args[0])
			}
			if err != nil {
				panic(err)
			}
			fmt.Print(str)
		},
	}
	prettyPrint bool
)

func init() {
	catFileCmd.Flags().
		BoolVarP(&prettyPrint, "pretty-print", "p", false, "Pretty-print the contents based on the type of the object")
}
