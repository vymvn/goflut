package utils

import (
	"image"
	"math"
	"net"
	"sync"

	vidio "github.com/AlexEidt/Vidio"
)

func checkCollision(imgWidth, imgHeight, x, y int, xvel, yvel *int) {

    if x + imgWidth >= canvasSize.width || x + imgWidth < 0 || x >= canvasSize.width || x <= 0 {
        *xvel *= -1
    }
    if y + imgHeight >= canvasSize.height || y + imgHeight <= 0 || y >= canvasSize.height || y <= 0 {
        *yvel *= -1
    }
}

func BouncyDrawVideo(videoPath string, startX, startY int, center bool, threads int, conn net.Conn) {
    xvel := 3
    yvel := 5
    getCanvasSize(&canvasSize, conn)
    for true {
        video, _ := vidio.NewVideo(videoPath)

        img := image.NewRGBA(image.Rect(0, 0, video.Width(), video.Height()))
        // img := image.NewRGBA(image.Rect(xvel, yvel, video.Width()+xvel*2, video.Height()+yvel*2))
        video.SetFrameBuffer(img.Pix)

        // var chunks []*chunk = makeChunks(threads, chunkWidth, scaledHeight) 

        width, height := getImageSize(img, conn)

        chunkWidth := width / threads
        var chunks []*chunk = newMakeChunks(threads, chunkWidth, height)


        if (center == true) {
            startX = (canvasSize.width / 2) - (width / 2)
            startY = (canvasSize.height / 2) - (height / 2)
        }

        frame := 0

        var wg sync.WaitGroup
        for video.Read() { 
            wg.Add(1)
            go newDrawChunk(chunks[0], img, startX, startY, &wg, conn)

            startX += xvel
            startY += yvel

            checkCollision(width, height, startX, startY, &xvel, &yvel)

            wg.Wait()
            // drawFrame(img, chunks, threads, startX, startY, size, center, conn)
            frame++
        }
    }
}


func BouncingVideo(videoPath string, startX, startY int, center bool, threads int, conn net.Conn) error {

    xvel := 5
    yvel := 5
    getCanvasSize(&canvasSize, conn)
    for true {
        video, err := vidio.NewVideo(videoPath)
        if err != nil {
            return err
        }

        img := image.NewRGBA(image.Rect(0, 0, video.Width(), video.Height()))
        video.SetFrameBuffer(img.Pix)

        width, height := getImageSize(img, conn)

        chunkWidth := width / threads
        var chunks []*chunk = newMakeChunks(threads, chunkWidth, height)

        if (center == true) {
            startX = (canvasSize.width / 2) - (width / 2)
            startY = (canvasSize.height / 2) - (height / 2)
        }

        var wg sync.WaitGroup
        frame := 0
        for video.Read() {

            checkCollision(width, height, startX, startY, &xvel, &yvel)
            startX += xvel
            startY += yvel

            for i := 0; i < threads; i++ {
                wg.Add(1)
                go newDrawChunk(chunks[i], img, startX, startY, &wg, conn)
            }
            wg.Wait()

            frame++
        }
    }

    return nil
}

func BouncingImage(img image.Image, x, y, xvel, yvel int, size float64, conn net.Conn) error {

    imgWidth, imgHeight, _ := getScaledImageSize(img, size, conn)

    drawCounter := 0
    for true {
        x += xvel
        y += yvel

        checkCollision(imgWidth, imgHeight, x, y, &xvel, &yvel)

        if (drawCounter == 10) {
            DrawImage(img, x, y, size, false, conn)
            drawCounter = 0
        }
        drawCounter++
    }
    return nil
}
