package cmd

import (
	"context"

	"os"

	"github.com/spf13/cobra"
)


var (
    host   string
    port   int
    startX int
    startY int
    center    bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goflut",
	Short: "A humble pixelflut client",
    // Run: runRoot,
}

var mainContext context.Context
func Execute() {
	// var cancel context.CancelFunc
	// mainContext, cancel = context.WithCancel(context.Background())
	// defer cancel()
	//
	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt)
	// defer func() {
	// 	signal.Stop(signalChan)
	// 	cancel()
	// }()
	// go func() {
	// 	select {
	// 	case <-signalChan:
	// 		// caught CTRL+C
	// 		fmt.Println("\n[!] Keyboard interrupt detected, terminating.")
	// 		cancel()
	// 	case <-mainContext.Done():
	// 	}
	// }()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
    rootCmd.PersistentFlags().StringVarP(&host, "host", "H", "", "The pixelflut server hostname or ip.")
    rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 0, "You know what to put here")
    rootCmd.PersistentFlags().IntVarP(&startX, "x-offset", "x", 0, "X axis offset.")
    rootCmd.PersistentFlags().IntVarP(&startY, "y-offset", "y", 0, "Y axis offset.")
    rootCmd.PersistentFlags().BoolVar(&center, "center", false, "Centers the drawing.")

    rootCmd.MarkPersistentFlagRequired("host")
    rootCmd.MarkPersistentFlagRequired("host")

}

// func runRoot(cmd *cobra.Command, args []string) {
//
//
// }

