package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vymvn/goflut/utils"
)


// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goflut",
	Short: "A humble pixelflut client",
    // Run: runRoot,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
    rootCmd.PersistentFlags().StringP("host", "H", "", "The pixelflut server hostname or ip.")
    rootCmd.PersistentFlags().IntP("port", "p", 0, "Server port.")
    rootCmd.PersistentFlags().IntP("x-offset", "x", 0, "X axis offset.")
    rootCmd.PersistentFlags().IntP("y-offset", "y", 0, "Y axis offset.")
    // rootCmd.PersistentFlags().Bool("center", false, "Centers the drawing.")
    rootCmd.PersistentFlags().Bool("loop", false, "Loops duh.")
    rootCmd.PersistentFlags().Int("threads", 1, "Number of threads to use.")

    rootCmd.MarkPersistentFlagRequired("host")
    rootCmd.MarkPersistentFlagRequired("host")

}

// func runRoot(cmd *cobra.Command, args []string) {
//
//
// }

func parseGlobalOptions() (*utils.GlobalOptions, error) {
	globalOpts := utils.NewGlobalOptions()

	threads, err := rootCmd.Flags().GetInt("threads")
	if err != nil {
		return nil, fmt.Errorf("invalid value for threads: %w", err)
	}

	if threads <= 0 {
		return nil, fmt.Errorf("threads must be bigger than 0")
	}
	globalOpts.Threads = threads

	globalOpts.Host, err = rootCmd.Flags().GetString("host")
	if err != nil {
		return nil, fmt.Errorf("invalid value for wordlist: %w", err)
	}

	port, err := rootCmd.Flags().GetInt("port")
	if err != nil {
		return nil, fmt.Errorf("invalid value for port: %w", err)
	}

	if port < 0 {
		return nil, fmt.Errorf("wordlist-offset must be bigger or equal to 0")
	}
	globalOpts.Port = port

	xOffset, err := rootCmd.Flags().GetInt("x-offset")
	if err != nil {
		return nil, fmt.Errorf("invalid value for x-offset: %w", err)
	}

	if xOffset < 0 {
		return nil, fmt.Errorf("x-offset must be bigger or equal to 0")
	}
	globalOpts.StartX = xOffset

	yOffset, err := rootCmd.Flags().GetInt("y-offset")
	if err != nil {
		return nil, fmt.Errorf("invalid value for y-offset: %w", err)
	}

	if xOffset < 0 {
		return nil, fmt.Errorf("y-offset must be bigger or equal to 0")
	}
	globalOpts.StartY = yOffset

	return globalOpts, nil
}
