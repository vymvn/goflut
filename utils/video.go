package utils

import (
	"image"
	"sync"
	"time"

	"github.com/AlexEidt/Vidio"
)

func InitVideo(videoPath string) (*vidio.Video, *image.RGBA) {

    video, _ := vidio.NewVideo(videoPath)

    frameBuffer := image.NewRGBA(image.Rect(0, 0, video.Width(), video.Height()))
    video.SetFrameBuffer(frameBuffer.Pix)

    return video, frameBuffer
}


func DrawVideo(video *vidio.Video, frameBuffer *image.RGBA, chunks []*VideoChunk, globalOpts *GlobalOptions, videoOpts *VideoOptions) {

    var wg sync.WaitGroup
    fps := video.FPS()
    totalFrames := video.Frames()

    for i := 0; i < len(chunks); i++ {
        chunks[i].currFrameBuffer = frameBuffer
    }

    // if (center == true) {
    //     startX = (canvasSize.width / 2) - (width / 2)
    //     startY = (canvasSize.height / 2) - (height / 2)
    // }

    if globalOpts.Loop {

        for true {

            for currFrame := 0; currFrame < totalFrames; currFrame++ {

                video.ReadFrame(currFrame)

                for i := 0; i < len(chunks); i++ {
                    wg.Add(1)
                    go DrawVideoChunk(chunks[i], globalOpts.StartX, globalOpts.StartY, &wg)
                }

                wg.Wait()
                time.Sleep(time.Second / time.Duration(fps))
            }

        }

    } else {

        for video.Read() {

            for i := 0; i < len(chunks); i++ {
                wg.Add(1)
                go DrawVideoChunk(chunks[i], globalOpts.StartX, globalOpts.StartY, &wg)
            }

            wg.Wait()
            time.Sleep(time.Second / time.Duration(fps))
        }
    }
}
