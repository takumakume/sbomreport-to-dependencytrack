package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/takumakume/sbomreport-to-dependencytrack/version"
)

var versiomCmd = &cobra.Command{
	Use:   "version",
	Short: "",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(version.Version)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versiomCmd)
}
