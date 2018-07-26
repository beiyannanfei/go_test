package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/beiyannanfei/go_test/0_auth_test/server"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
		and usage of using your command. For example:

		Cobra is a CLI library for Go that empowers applications.
		This application is a tool to generate the needed files
		to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		server.Run()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().BoolP("graceful", "g", false, "listen on fd open 3 (internal use only)")
	viper.BindPFlag("graceful", startCmd.Flags().Lookup("graceful"))
}
