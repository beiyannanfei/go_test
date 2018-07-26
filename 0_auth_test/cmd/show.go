package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"strings"
)

var echoTimes int

var showCmd = &cobra.Command{
	Use:   "show [string to show]",
	Short: "Echo anything to the screen",
	Long: `echo is for echoing anything back.
		   Echo works a lot like print, except it has a child command.
		`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("show: " + strings.Join(args, ", "))
	},
}

var timesCmd = &cobra.Command{
	Use:   "times [# times] [string ro echo]",
	Short: "Echo anything to the screen more times",
	Long: `echo things multiple times back to the user by providing
		a count and a string.`,
	Run: func(cmd *cobra.Command, args []string) {
		for i := 0; i < echoTimes; i++ {
			fmt.Println("Echo: " + strings.Join(args, ", "))
		}
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
	showCmd.AddCommand(timesCmd)

	timesCmd.Flags().IntVarP(&echoTimes, "times", "t", 1, "times to echo the intput")
}
