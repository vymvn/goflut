package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vymvn/goflut/utils"
)

var videoCmd = &cobra.Command{
	Use:   "video",
	Short: "Video streaming mode.",
	Run: runVideo,
}

func init() {
	rootCmd.AddCommand(videoCmd)

	videoCmd.Flags().StringP("video", "v", "", "Path of video or gif. (Or anything that ffmpeg supports lol)")
    // videoCmd.Flags().Float64VarP(&videoSize, "size", "s", 1, "Size of the image where 1 is the original size.")
    videoCmd.Flags().Bool("bounce", false, "Bounce around (best used with a smaller video/gif)")
    videoCmd.Flags().Bool("center", false, "Centers the video on the canvas.")
    videoCmd.Flags().Float64("x-vel", 70, "The velocity on the X-Axis. (only for bounce mode)")
    videoCmd.Flags().Float64("y-vel", 70, "The velocity on the Y-Axis. (only for bounce mode)")

    videoCmd.MarkFlagRequired("video")
}

func runVideo(cmd *cobra.Command, args []string) {

    // Parsing flags
    globalOpts, videoOpts, err := parseVideoOptions()
    if err != nil {
        fmt.Fprintln(os.Stderr, "error on parsing arguments: %w", err)
        os.Exit(1)
    }

    video, frameBuffer := utils.InitVideo(videoOpts.Path)

    if videoOpts.Bounce {

        chunks := utils.MakeVideoChunks(frameBuffer, globalOpts, videoOpts)
        utils.BouncyDrawVideo(video, frameBuffer, chunks, globalOpts, videoOpts)

        // utils.BouncyDrawVideo(videoPath, startX, startY, center, videoThreads, conn)
    } else {

        chunks := utils.MakeVideoChunks(frameBuffer, globalOpts, videoOpts)
        utils.DrawVideo(video, frameBuffer, chunks, globalOpts, videoOpts)
    }

}


func parseVideoOptions() (*utils.GlobalOptions, *utils.VideoOptions, error) {

    globalOpts, err := parseGlobalOptions()
    if err != nil {
        return nil, nil, err
    }

    videoOpts := utils.NewVideoOptions()

    videoOpts.Bounce, err = imageCmd.Flags().GetBool("bounce")
    if err != nil {
        return nil, nil, fmt.Errorf("could not set bounce flag: %w", err)
    }

    videoOpts.VelocityX, err = imageCmd.Flags().GetFloat64("x-vel")
    if err != nil {
        return nil, nil, fmt.Errorf("invalid value for image x-velocity: %w", err)
    }

    videoOpts.VelocityY, err = imageCmd.Flags().GetFloat64("y-vel")
    if err != nil {
        return nil, nil, fmt.Errorf("invalid value for image y-velocity: %w", err)
    }

    return globalOpts, videoOpts, nil
}
