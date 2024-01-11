package cmd

import (
	"fmt"
	"net"
	"os"

	"github.com/spf13/cobra"
	"github.com/vymvn/goflut/utils"
)

var textCmd = &cobra.Command{
	Use:   "text",
	Short: "Text rendering mode",
	Run: runText,
}

var (
    text      string
    fontSize  float64
)

func init() {
	textCmd.Flags().StringVarP(&text, "text", "t", "", "Text to be rendered.")
	textCmd.Flags().Float64VarP(&fontSize, "font-size", "s", 12, "Font size of text.")

    textCmd.MarkFlagRequired("text")

	rootCmd.AddCommand(textCmd)
}

func runText(cmd *cobra.Command, args []string) {

    connString := fmt.Sprintf("%s:%d", host, port)
    conn, err := net.Dial("tcp", connString)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Could not connect to \"" + connString + "\":\n", err)
        os.Exit(1)
    }
    defer conn.Close()

    err = utils.DrawText(text, startX, startY, fontSize, center, conn)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Could not draw text:\n", err)
    }

}
