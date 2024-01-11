package cmd

import (

	"github.com/spf13/cobra"
)

// bounceCmd represents the bounce command
var bounceCmd = &cobra.Command{
	Use:   "bounce",
	Short: "Bouncing mode",
	Run: func(cmd *cobra.Command, args []string) {
		bounce = true
	},
}

var (
     bounce bool
)

func init() {

	bounceCmd.Flags().Float64("x-vel", 1, "The velocity on the X-Axis.")
	bounceCmd.Flags().Float64("y-vel", 2, "The velocity on the Y-Axis.")

	rootCmd.AddCommand(bounceCmd)
}
