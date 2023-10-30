package cmd

import (
	"fmt"
	"os"

	"github.com/magnickolas/gitok/config/parser"
	"github.com/spf13/cobra"
)

var (
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "Manage config",
		Run: func(cmd *cobra.Command, args []string) {
			if list {
				r, err := os.Open(".git/config")
				if err != nil {
					fatalf("cannot open %v: %v", ".git/config", err)
				}
				p, err := parser.NewParser(r)
				if err != nil {
					fatalf("cannot init parser %v: %v", ".git/config", err)
				}
				kvs, err := p.Parse()
				if err != nil {
					fatalf("cannot parse %v: %v", ".git/config", err)
				}
				for _, kv := range kvs {
					fmt.Printf("%s=%s\n", kv.Key, kv.Value)
				}
			}
		},
	}
	list bool
)

func init() {
	configCmd.Flags().
		BoolVarP(&list, "list", "l", false, "List config key values")
}
