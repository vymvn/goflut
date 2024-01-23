package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vymvn/goflut/utils"
)

var textCmd = &cobra.Command{
	Use:   "text",
	Short: "Text rendering mode",
	Run: runText,
}

func init() {
	textCmd.Flags().StringP("text", "t", "", "Text to be rendered.")
	textCmd.Flags().Float64P("font-size", "s", 12, "Font size of text.")
	textCmd.Flags().StringP("font", "f", "fonts/Lato-Regular.ttf", "Font size of text.")

    textCmd.MarkFlagRequired("text")

	rootCmd.AddCommand(textCmd)
}

func runText(cmd *cobra.Command, args []string) {

    // Parsing flags
    globalOpts, textOpts, err := parseTextOptions()
    if err != nil {
        fmt.Errorf("error on parsing arguments: %w", err)
    }

    if globalOpts.Loop == true {

        for true {

            err = utils.DrawText(globalOpts, textOpts)
            if err != nil {
                fmt.Fprintln(os.Stderr, "Could not draw text:\n", err)
            }

        }

    } else {

        err = utils.DrawText(globalOpts, textOpts)
        if err != nil {
            fmt.Fprintln(os.Stderr, "Could not draw text:\n", err)
        }

    }

}

func parseTextOptions() (*utils.GlobalOptions, *utils.TextOptions, error) {

    globalOpts, err := parseGlobalOptions()
    if err != nil {
        return nil, nil, err
    }

    textOpts := utils.NewTextOptions()

    textOpts.Text, err  = imageCmd.Flags().GetString("text")
    if err != nil {
        return nil, nil, fmt.Errorf("invalid value for text: %w", err)
    }

    textOpts.FontPath, err  = imageCmd.Flags().GetString("font")
    if err != nil {
        return nil, nil, fmt.Errorf("invalid value for text: %w", err)
    }

    textOpts.FontSize, err  = imageCmd.Flags().GetFloat64("font-size")
    if err != nil {
        return nil, nil, fmt.Errorf("invalid value for text: %w", err)
    }

    return globalOpts, textOpts, nil
}
