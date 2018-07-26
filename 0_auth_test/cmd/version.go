package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of loveauth",
	Long:  `All software has versions. This is loveauth's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("loveauth v1.0.66 -- online(3e31d20)")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
