package utils

import (
	"fmt"
	"image"
	"sync"
	"time"

	vidio "github.com/AlexEidt/Vidio"
)

func checkCollision(imgWidth, imgHeight int, x, y int, xvel, yvel *float64) {

    if x + imgWidth >= canvasSize.width || x + imgWidth < 0 || x >= canvasSize.width || x <= 0 {
        *xvel *= -1
    }
    if y + imgHeight >= canvasSize.height || y + imgHeight <= 0 || y >= canvasSize.height || y <= 0 {
        *yvel *= -1
    }
}

func BouncyDrawVideo(video *vidio.Video, frameBuffer *image.RGBA, chunks []*VideoChunk, globalOpts *GlobalOptions, videoOpts *VideoOptions) error {

    connString := fmt.Sprintf("%s:%d", globalOpts.Host, globalOpts.Port)
    err, conn := makeConnection(connString)
    if err != nil {
        return err
    }
    getCanvasSize(&canvasSize, conn)
    conn.Close()

    fps := video.FPS()
    // fps := 60
    totalFrames := video.Frames()

    frameWidth  := frameBuffer.Bounds().Max.X
    frameHeight := frameBuffer.Bounds().Max.Y

    xvel := 50.0
    yvel := 50.0

    // if (center == true) {
    //     startX = (canvasSize.width / 2) - (width / 2)
    //     startY = (canvasSize.height / 2) - (height / 2)
    // }

    for i := 0; i < len(chunks); i++ {
        chunks[i].currFrameBuffer = frameBuffer
    }

    var wg sync.WaitGroup
    for true {

        lastTime := time.Now()
        // for video.Read() { 
        for currFrame := 0; currFrame < totalFrames; currFrame++ {

            video.ReadFrame(currFrame)

            // Calculate delta time
            currentTime := time.Now()
            dt := currentTime.Sub(lastTime).Seconds()
            lastTime = currentTime

            for i := 0; i < len(chunks); i++ {
                wg.Add(1)
                // chunks[i].currFrameBuffer = frameBuffer
                go DrawVideoChunkFull(chunks[i], globalOpts.StartX, globalOpts.StartY, &wg)
            }

            globalOpts.StartX += int(xvel * dt)
            globalOpts.StartY += int(yvel * dt)
            // xOffset += int(xvel)
            // yOffset += int(yvel)

            checkCollision(frameWidth, frameHeight, globalOpts.StartX, globalOpts.StartY, &xvel, &yvel)

            wg.Wait()
            // drawFrame(img, chunks, threads, startX, startY, size, center, conn)
            time.Sleep(time.Second / time.Duration(fps))
        }

    }

    return nil
}

func BouncingImage(chunks []*ImageChunk, globalOpts *GlobalOptions, imageOpts *ImageOptions) error {

    connString := fmt.Sprintf("%s:%d", globalOpts.Host, globalOpts.Port)
    err, conn := makeConnection(connString)
    if err != nil {
        return err
    }

    getCanvasSize(&canvasSize, conn)

    imgWidth, imgHeight, _ := getScaledImageSize(chunks[0].img, imageOpts.Scale, conn)
    conn.Close()

    var wg sync.WaitGroup
    lastTime := time.Now()
    for true {

        // drawCounter := 0

        // Calculate delta time
        currentTime := time.Now()
        dt := currentTime.Sub(lastTime).Seconds()
        lastTime = currentTime

        globalOpts.StartX += int(imageOpts.VelocityX * dt)
        globalOpts.StartY += int(imageOpts.VelocityY * dt)
        // xOffset += int(xvel)
        // yOffset += int(yvel)


        // if (drawCounter == 10) {
        // drawCounter = 0
        // }
        // drawCounter++
        checkCollision(imgWidth, imgHeight, globalOpts.StartX, globalOpts.StartY, &imageOpts.VelocityX, &imageOpts.VelocityY)
        for i := 0; i < len(chunks); i++ {
            wg.Add(1)
            go expDrawImageChunk(chunks[i], globalOpts.StartX, globalOpts.StartY, &wg)
        }

        wg.Wait()

    }
    return nil
}
