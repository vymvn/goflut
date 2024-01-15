package utils

import (
	"image"
	"net"
	"sync"

	"github.com/AlexEidt/Vidio"
)


// func drawFrame(img image.Image, chunks []*chunk, threads int, startX int, startY int, size float64, center bool, conn net.Conn) error {
//
//     // scaledWidth, scaledHeight, scale := getScaledImageSize(img, size, conn)
//
//     if (center == true) {
//         startX = (canvasSize.width / 2) - (scaledWidth / 2)
//         startY = (canvasSize.height / 2) - (scaledHeight / 2)
//     }
//
//     var wg sync.WaitGroup
//     for i := 0; i < threads; i++ {
//         wg.Add(1)
//         go drawChunk(chunks[i], img, startX, startY, &wg, conn)
//     }
//
//     wg.Wait()
//
//
//     return nil
// }


func NewDrawVideo(videoPath string, startX, startY int, center bool, threads int, conn net.Conn) {
    video, _ := vidio.NewVideo(videoPath)

    img := image.NewRGBA(image.Rect(0, 0, video.Width(), video.Height()))
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

        for i := 0; i < threads; i++ {
            wg.Add(1)
            go newDrawChunk(chunks[i], img, startX, startY, &wg, conn)
        }

        wg.Wait()
        // drawFrame(img, chunks, threads, startX, startY, size, center, conn)
        frame++
    }
}

func DrawVideo(videoPath string, startX, startY int, size float64, threads int, center bool, conn net.Conn) {
    video, _ := vidio.NewVideo(videoPath)

    img := image.NewRGBA(image.Rect(0, 0, video.Width(), video.Height()))
    video.SetFrameBuffer(img.Pix)

    scaledWidth, scaledHeight, scale := getScaledImageSize(img, size, conn)

    chunkWidth := scaledWidth / threads
    var chunks []*chunk = makeChunks(threads, chunkWidth, scaledHeight, scale)

    if (center == true) {
        startX = (canvasSize.width / 2) - (scaledWidth / 2)
        startY = (canvasSize.height / 2) - (scaledHeight / 2)
    }

    var wg sync.WaitGroup
    frame := 0

    for video.Read() {
        for i := 0; i < threads; i++ {
            wg.Add(1)
            go drawChunk(chunks[i], img, startX, startY, &wg, conn)
        }

        wg.Wait()

        frame++
    }
}
