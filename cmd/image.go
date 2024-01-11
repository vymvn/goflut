package cmd

import (
	"fmt"
	"net"
	"os"

	"github.com/spf13/cobra"
	"github.com/vymvn/goflut/utils"
)

var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Draws the image with the passed path and cords.",
	Run: runImage,
}

var (
    imagePath string
    size      float64
    center    bool
)

func init() {
	// bounceCmd.AddCommand(imageCmd)

	imageCmd.Flags().Float64VarP(&size, "size", "s", 1, "Size of the image where 1 is the original size.")
	// imageCmd.Flags().BoolP("scale-to-fit", "S", true, "Scale to fit the image to the canvas maintaing aspect ratio.")
	imageCmd.Flags().StringVarP(&imagePath, "image", "i", "", "Path to the image to draw. (required)")
    imageCmd.MarkFlagRequired("image")
	imageCmd.Flags().IntVar(&startX, "x", 0, "Starting X")
	imageCmd.Flags().IntVar(&startY, "y", 0, "Starting Y")
	imageCmd.Flags().BoolVar(&center, "center", false, "Center image on canvas")

	rootCmd.AddCommand(imageCmd)
    // if (imageCmd.Parent().Use == "bounce") {
    //     bounce = true
    // }

}

func runImage(cmd *cobra.Command, args []string) {

    connString := fmt.Sprintf("%s:%d", host, port)
    conn, err := net.Dial("tcp", connString)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Could not connect to \"" + connString + "\":\n", err)
        os.Exit(1)
    }
    defer conn.Close()

    err = utils.DrawImageFromPath(imagePath, startX, startY, size, center, conn)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Could not draw image:\n", err)
    }
}
