package cmd

import (
	"fmt"
	"net"
	"os"

	"github.com/spf13/cobra"
	"github.com/vymvn/goflut/utils"
)

var wipeCmd = &cobra.Command{
	Use:   "wipe",
	Short: "Wipes the canvas.",
    Run: func(cmd *cobra.Command, args []string) {

        connString := fmt.Sprintf("%s:%d", host, port)
        conn, err := net.Dial("tcp", connString)
        if err != nil {
            fmt.Fprintln(os.Stderr, "Could not connect to \"" + connString + "\":\n", err)
            os.Exit(1)
        }
        defer conn.Close()

        utils.WipeCanvas(conn)
    },
}

func init() {
	rootCmd.AddCommand(wipeCmd)
}
