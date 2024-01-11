package cmd

import (
    "fmt"
    "image"
    "net"
    "os"

    "github.com/spf13/cobra"
    "github.com/vymvn/goflut/utils"
)

var imageCmd = &cobra.Command{
    Use:   "image",
    Short: "Image drawing mode.",
    Run: runImage,
}

var (
    bounce    bool
    xVel   int
    yVel   int
    imagePath string
    size      float64
    center    bool
)

func init() {

    imageCmd.Flags().Float64VarP(&size, "size", "s", 1, "Size of the image where 1 is the original size.")
    // imageCmd.Flags().BoolP("scale-to-fit", "S", true, "Scale to fit the image to the canvas maintaing aspect ratio.")
    imageCmd.Flags().StringVarP(&imagePath, "image", "i", "", "Path to the image to draw. (required)")
    imageCmd.MarkFlagRequired("image")
    imageCmd.Flags().IntVar(&startX, "x", 0, "Starting X")
    imageCmd.Flags().IntVar(&startY, "y", 0, "Starting Y")
    imageCmd.Flags().BoolVar(&center, "center", false, "Center image on canvas")
    imageCmd.Flags().BoolVar(&bounce, "bounce", false, "Bounce around (best used with a smaller picture)")
	imageCmd.Flags().IntVar(&xVel, "x-vel", 1, "The velocity on the X-Axis. (only for bounce mode)")
	imageCmd.Flags().IntVar(&yVel, "y-vel", 2, "The velocity on the Y-Axis. (only for bounce mode)")

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

    f, err := os.Open(imagePath)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Failed to open image:\n", err)
    }
    defer f.Close()

    img, _, err := image.Decode(f)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Failed to decode image:\n", err)
    }

    if (bounce == true) {
        err := utils.BouncingImage(img, startX, startY, xVel, yVel, size, conn)
        if err != nil {
            fmt.Fprintln(os.Stderr, "Failed to bounce image:\n", err)
        }

    } else {

        err = utils.DrawImage(img, startX, startY, size, center, conn)
        if err != nil {
            fmt.Fprintln(os.Stderr, "Could not draw image:\n", err)
        }

    }

}
