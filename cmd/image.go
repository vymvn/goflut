package cmd

import (
    "fmt"
    "image"
    "os"

    "github.com/spf13/cobra"
    "github.com/vymvn/goflut/utils"
)

type ImageOptions struct {
    Path      string
    Scale     float64
    Bounce    bool
    Center    bool
    VelocityX float64
    VelocityY float64

}

var imageCmd *cobra.Command

func runImage(cmd *cobra.Command, args []string) {

    // Parsing flags
    globalOpts, imageOpts, err := parseImageOptions()
    if err != nil {
        fmt.Errorf("error on parsing arguments: %w", err)
    }

    // Opening image
    f, err := os.Open(imageOpts.Path)
    if err != nil {
        fmt.Errorf("error opening image: %w", err)
    }
    defer f.Close()

    // Decoding image into an image.Image struct
    img, _, err := image.Decode(f)
    if err != nil {
        fmt.Errorf("error decoding image: %w", err)
    }

    // Splitting image into chunks for threading
    err, chunks := utils.ExpMakeImageChunks(img, globalOpts, imageOpts)
    if err != nil {
        fmt.Errorf("error making image chunks: %w", err)
    }

    if (imageOpts.Bounce == true) {

        err = utils.BouncingImage(chunks, globalOpts, imageOpts)
        if err != nil {
            fmt.Fprintln(os.Stderr, "Failed to bounce image:\n", err)
        }

    } else if (globalOpts.Loop == true) {

        for true {
            err := utils.ExpDrawImageThreaded(chunks, globalOpts)
            // err = utils.DrawImageThreaded(img, startX, startY, imageSize, imageThreads, center, conn)
            if err != nil {
                fmt.Fprintln(os.Stderr, "Could not draw image:\n", err)
                os.Exit(1)
            }
        }

    } else {

        err = utils.ExpDrawImageThreaded(chunks, globalOpts)
        // err = utils.DrawImageThreaded(img, startX, startY, imageSize, imageThreads, center, conn)
        if err != nil {
            fmt.Fprintln(os.Stderr, "Could not draw image:\n", err)
            os.Exit(1)
        }
    }

}

func parseImageOptions() (*utils.GlobalOptions, *utils.ImageOptions, error) {
    globalOpts, err := parseGlobalOptions()
    if err != nil {
        return nil, nil, err
    }

    imageOpts := utils.NewImageOptions()

    imageOpts.Path, err  = imageCmd.Flags().GetString("path")
    if err != nil {
        return nil, nil, fmt.Errorf("invalid value for image path: %w", err)
    }

    imageOpts.Scale, err = imageCmd.Flags().GetFloat64("scale")
    if err != nil {
        return nil, nil, fmt.Errorf("invalid value for image scale: %w", err)
    }

    imageOpts.Bounce, err = imageCmd.Flags().GetBool("bounce")
    if err != nil {
        return nil, nil, fmt.Errorf("could not set bounce flag: %w", err)
    }

    imageOpts.VelocityX, err = imageCmd.Flags().GetFloat64("x-vel")
    if err != nil {
        return nil, nil, fmt.Errorf("invalid value for image x-velocity: %w", err)
    }

    imageOpts.VelocityY, err = imageCmd.Flags().GetFloat64("y-vel")
    if err != nil {
        return nil, nil, fmt.Errorf("invalid value for image y-velocity: %w", err)
    }

    return globalOpts, imageOpts, nil
}

func init() {

    imageCmd = &cobra.Command{
        Use:   "image",
        Short: "Image drawing mode.",
        Run: runImage,
    }

    imageCmd.Flags().Float64P("scale", "s", 1, "Scale of the image where 1 is the original size.")
    // imageCmd.Flags().BoolP("scale-to-fit", "S", true, "Scale to fit the image to the canvas maintaing aspect ratio.")
    imageCmd.Flags().StringP("image", "i", "", "Path to the image to draw. (required)")
    imageCmd.MarkFlagRequired("image")
    // imageCmd.Flags().IntVar(&startX, "x", 0, "Starting X")
    // imageCmd.Flags().IntVar(&startY, "y", 0, "Starting Y")
    // imageCmd.Flags().BoolVar(&center, "center", false, "Center image on canvas")
    imageCmd.Flags().Bool("bounce", false, "Bounce around (best used with a smaller picture)")
    imageCmd.Flags().Bool("center", false, "Centers the image on the canvas.")
    imageCmd.Flags().Float64("x-vel", 70, "The velocity on the X-Axis. (only for bounce mode)")
    imageCmd.Flags().Float64("y-vel", 70, "The velocity on the Y-Axis. (only for bounce mode)")

    rootCmd.AddCommand(imageCmd)

}
