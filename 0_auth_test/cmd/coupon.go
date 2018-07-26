package cmd

import (
	"github.com/spf13/viper"
	"github.com/spf13/cobra"
	"fmt"
)

var couponCmd = &cobra.Command{
	Use:   "coupon",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
		and usage of using your command. For example:

		Cobra is a CLI library for Go that empowers applications.
		This application is a tool to generate the needed files
		to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		gen := viper.GetBool("generate")
		val := viper.GetBool("validate")
		ben := viper.GetBool("benchmark")
		fmt.Printf("run coupon: gen: %#v, val: %#v, ben: %#v\n", gen, val, ben)
	},
}

func init() {
	rootCmd.AddCommand(couponCmd)

	couponCmd.Flags().BoolP("generate", "g", false, "Generate coupons")
	couponCmd.Flags().BoolP("validate", "v", false, "Validate coupons")
	couponCmd.Flags().BoolP("benchmark", "b", false, "BenchMark coupons")

	viper.BindPFlag("generate", couponCmd.Flags().Lookup("generate"))
	viper.BindPFlag("validate", couponCmd.Flags().Lookup("validate"))
	viper.BindPFlag("benchmark", couponCmd.Flags().Lookup("benchmark"))
}
