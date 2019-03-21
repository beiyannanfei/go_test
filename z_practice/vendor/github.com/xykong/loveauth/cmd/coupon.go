// Copyright © 2018 xykong <xy.kong@gmail.com>

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xykong/loveauth/coupon"
)

// couponCmd represents the coupon command
var couponCmd = &cobra.Command{
	Use:   "coupon",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("coupon called")

		if viper.GetBool("benchmark") {
			coupon.Benchmark()
		}
	},
}

func init() {
	rootCmd.AddCommand(couponCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// couponCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// couponCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	couponCmd.Flags().BoolP("generate", "g", false, "Generate coupons")
	couponCmd.Flags().BoolP("validate", "v", false, "Validate coupons")
	couponCmd.Flags().BoolP("benchmark", "b", false, "Benchmark coupons")

	viper.BindPFlag("generate", couponCmd.Flags().Lookup("generate"))
	viper.BindPFlag("validate", couponCmd.Flags().Lookup("validate"))
	viper.BindPFlag("benchmark", couponCmd.Flags().Lookup("benchmark"))
}
