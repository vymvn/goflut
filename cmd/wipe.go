package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vymvn/goflut/utils"
)

var wipeCmd = &cobra.Command {
	Use:   "wipe",
	Short: "Wipes the canvas.",
    Run: func(cmd *cobra.Command, args []string) {

        globalOpts, err := parseGlobalOptions()
        if err != nil {
            fmt.Fprintln(os.Stderr, "error on parsing arguments: %w", err)
            os.Exit(1)
        }
        connString := fmt.Sprintf("%s:%d", globalOpts.Host, globalOpts.Port)

        err = utils.WipeCanvas(connString)
        // err := utils.Noise(startX, startY, connString)
        if err != nil {
            fmt.Fprintln(os.Stderr, "Could not wipe canvas:\n", err)
            os.Exit(1)
        }
    },
}

func init() {
	rootCmd.AddCommand(wipeCmd)
}
