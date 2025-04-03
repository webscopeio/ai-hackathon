package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/webscopeio/ai-hackathon/internal/repository/analyze"
)

var rootCmd = &cobra.Command{
	Use:   "testbuddy",
	Short: "TestBuddy CLI",
	Run: func(cmd *cobra.Command, args []string) {
		err := analyze.Analyze()
		if err != nil {
			fmt.Println(err)
		}
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
