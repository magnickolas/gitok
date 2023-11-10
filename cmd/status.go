package cmd

import (
	"fmt"
	"os"

	"github.com/magnickolas/gitok/index/parser"
	"github.com/spf13/cobra"
)

var (
	statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Show current index",
		Run: func(cmd *cobra.Command, args []string) {
			r, err := os.Open(".git/index")
			if err != nil {
				fatalf("cannot open %v: %v", args[0], err)
			}
			p, err := parser.NewParser(r)
			if err != nil {
				fatalln("failed to create a parser")
			}
			index, err := p.Parse()
			if err != nil {
				fatalf("failed to parse index: %v\n", err)
			}
			for _, entry := range index.Entries {
				fmt.Println(entry.Name)
			}
		},
	}
)

func init() {
}
