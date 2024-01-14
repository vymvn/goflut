package cmd

import (
	"fmt"
	"net"
	"os"

	"github.com/spf13/cobra"
	"github.com/vymvn/goflut/utils"
)

var videoCmd = &cobra.Command{
	Use:   "video",
	Short: "Video streaming mode.",
	Run: runVideo,
}

var (
    videoPath    string
    videoLoop    bool
    // videoSize    float64
    videoThreads int
)

func init() {
	rootCmd.AddCommand(videoCmd)

	videoCmd.Flags().StringVarP(&videoPath, "video", "v", "", "Help message for toggle")
    // videoCmd.Flags().Float64VarP(&videoSize, "size", "s", 1, "Size of the image where 1 is the original size.")
    videoCmd.Flags().BoolVar(&videoLoop, "loop", false, "Keeps drawing in a loop.")
    videoCmd.Flags().IntVar(&videoThreads, "threads", 1, "Number of threads.")
    videoCmd.MarkFlagRequired("video")
}

func runVideo(cmd *cobra.Command, args []string) {


    connString := fmt.Sprintf("%s:%d", host, port)
    conn, err := net.Dial("tcp", connString)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Could not connect to \"" + connString + "\":\n", err)
        os.Exit(1)
    }
    defer conn.Close()

    if videoLoop {

        for true {
            // utils.DrawVideo(videoPath, startX, startY, videoSize, videoThreads, center, conn)
            utils.NewDrawVideo(videoPath, startX, startY, center, videoThreads, conn)
        }
    } else {
        // utils.DrawVideo(videoPath, startX, startY, videoSize, videoThreads, center, conn)
        utils.NewDrawVideo(videoPath, startX, startY, center, videoThreads, conn)
    }

}
