package cmd

import (
	"fmt"
	"os"

	"github.com/magnickolas/gitok/gitok_hash"
	"github.com/spf13/cobra"
)

var (
	hashObjectCmd = &cobra.Command{
		Use:   "hash-object",
		Short: "Hash an object",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if !readFromStdin && len(args) == 0 {
				panic(fmt.Errorf("expected filename"))
			}
			var r *os.File
			if readFromStdin {
				r = os.Stdin
			} else {
				var err error
				r, err = os.Open(args[0])
				if err != nil {
					panic(err)
				}
			}
			key, err := gitok_hash.ProcessBlob(r, write)
			if err != nil {
				panic(err)
			}
			fmt.Println(key)
		},
	}
	write         bool
	readFromStdin bool
)

func init() {
	hashObjectCmd.Flags().
		BoolVarP(&write, "write", "w", false, "whether to write an object to storage")
	hashObjectCmd.Flags().
		BoolVar(&readFromStdin, "stdin", false, "read from stdin instead of file")
}
