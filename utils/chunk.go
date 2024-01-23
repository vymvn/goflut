package utils

import (
	"fmt"
	"image"
	"net"
	"os"
	"sync"

)


type VideoChunk struct {
    xPos   int
    width  int
    height int
    prevFrameBuffer *image.RGBA
    currFrameBuffer *image.RGBA
    conn   net.Conn
}

type ImageChunk struct {
    xPos   int
    width  int
    height int
    scale  float64
    img    image.Image
    buffer *image.RGBA
    conn   net.Conn
}


func MakeVideoChunks(frameBuffer *image.RGBA, globalOpts *GlobalOptions, videoOpts *VideoOptions) []*VideoChunk {

    chunks := make([]*VideoChunk, globalOpts.Threads)   // As many chunks as threads

    frameWidth  := frameBuffer.Bounds().Max.X
    frameHeight := frameBuffer.Bounds().Max.Y

    connString := fmt.Sprintf("%s:%d", globalOpts.Host, globalOpts.Port)

    // if (center == true) {
    //     startX = (canvasSize.width / 2) - (scaledWidth / 2)
    //     startY = (canvasSize.height / 2) - (scaledHeight / 2)
    // }

    // connString := fmt.Sprintf("localhost:1234")
    // connString := fmt.Sprintf("pixelflut.uwu.industries:1234")

    chunkWidth := frameWidth / globalOpts.Threads
    currIndex := 0
    for i := 0; i < len(chunks); i++{

        newConn, err := net.Dial("tcp", connString)
        if err != nil {
            fmt.Fprintln(os.Stderr, "Could not connect to \"" + connString + "\":\n", err)
            os.Exit(1)
        }
        initFrameBuffer := image.NewRGBA(image.Rect(0, 0, frameWidth, frameHeight))
        chunks[i] = &VideoChunk{
            xPos   : currIndex,
            width  : chunkWidth,
            height : frameHeight,
            prevFrameBuffer: initFrameBuffer,
            conn   : newConn,
        }

        currIndex += chunkWidth
    }

    return chunks
}

func ExpMakeImageChunks(img image.Image, globalOpts *GlobalOptions, imageOpts *ImageOptions) (error, []*ImageChunk) {

    chunks := make([]*ImageChunk, globalOpts.Threads)   // As many chunks as threads

    connString := fmt.Sprintf("%s:%d", globalOpts.Host, globalOpts.Port)
    err, conn := makeConnection(connString)
    if err != nil {
        return err, nil
    }
    scaledWidth, scaledHeight, scaleFactor := getScaledImageSize(img, imageOpts.Scale, conn)

    if (imageOpts.Center == true) {
        globalOpts.StartX = (canvasSize.width / 2) - (scaledWidth / 2)
        globalOpts.StartY = (canvasSize.height / 2) - (scaledHeight / 2)
    }

    chunkWidth := scaledWidth / globalOpts.Threads
    currIndex := 0
    for i := 0; i < len(chunks); i++{

        newConn, err := net.Dial("tcp", connString)
        if err != nil {
            fmt.Fprintln(os.Stderr, "Could not connect to \"" + connString + "\":\n", err)
            os.Exit(1)
        }
        chunkBuffer := image.NewRGBA(image.Rect(0, 0, chunkWidth, scaledHeight))
        chunks[i] = &ImageChunk{
            xPos   : currIndex,
            width  : chunkWidth,
            height : scaledHeight,
            scale  : scaleFactor,
            buffer : chunkBuffer,
            img    : img,
            conn   : newConn,
        }

        currIndex += chunkWidth
    }

    return nil, chunks
}


func DrawVideoChunkFull(chunk *VideoChunk, xOffset, yOffset int, wg *sync.WaitGroup) {

    defer wg.Done()

    for x := chunk.xPos; x < (chunk.xPos + chunk.width); x++ {
        for y := 0; y < chunk.height; y++ {

            currPixelColor := chunk.currFrameBuffer.At(x, y)
            r, g, b, a := currPixelColor.RGBA()
            WritePixel(x + xOffset, y + yOffset, int(r>>8), int(g>>8), int(b>>8), int(a>>8), chunk.conn)
        }
    }
}

func DrawVideoChunk(chunk *VideoChunk, xOffset, yOffset int, wg *sync.WaitGroup) {

    defer wg.Done()

    for x := chunk.xPos; x < (chunk.xPos + chunk.width); x++ {
        for y := 0; y < chunk.height; y++ {

            currPixelColor := chunk.currFrameBuffer.At(x, y)
            if (currPixelColor != chunk.prevFrameBuffer.At(x, y)) {
                r, g, b, a := currPixelColor.RGBA()
                WritePixel(x + xOffset, y + yOffset, int(r>>8), int(g>>8), int(b>>8), int(a>>8), chunk.conn)
                chunk.prevFrameBuffer.Set(x, y, currPixelColor)
            }
        }
    }
}

func expDrawImageChunk(chunk *ImageChunk, xOffset, yOffset int, wg *sync.WaitGroup) {

    defer wg.Done()

    for x := chunk.xPos; x < (chunk.xPos + chunk.width); x++ {
        for y := 0; y < chunk.height; y++ {
            scaledX := int(float64(x) / chunk.scale)
            scaledY := int(float64(y) / chunk.scale)

            currPixelColor := chunk.img.At(scaledX, scaledY)
            r, g, b, a := currPixelColor.RGBA()
            WritePixel(x + xOffset, y + yOffset, int(r>>8), int(g>>8), int(b>>8), int(a>>8), chunk.conn)
        }
    }
}

func newMakeChunks(threadsCount int, chunkWidth int, chunkHeight int) []*ImageChunk {

    chunks := make([]*ImageChunk, threadsCount)   // As many chunks as threads

    currIndex := 0
    for i := 0; i < len(chunks); i++{

        chunkBuffer := image.NewRGBA(image.Rect(0, 0, chunkWidth, chunkHeight))
        chunks[i] = &ImageChunk{
            xPos   : currIndex,
            width  : chunkWidth,
            height : chunkHeight,
            buffer : chunkBuffer,
        }

        currIndex += chunkWidth
    }

    return chunks
}

func makeChunks(threadsCount int, chunkWidth int, chunkHeight int, chunkScale float64) []*ImageChunk {

    chunks := make([]*ImageChunk, threadsCount)   // As many chunks as threads

    currIndex := 0
    for i := 0; i < len(chunks); i++{

        chunkBuffer := image.NewRGBA(image.Rect(0, 0, chunkWidth, chunkHeight))
        chunks[i] = &ImageChunk{
            xPos   : currIndex,
            width  : chunkWidth,
            height : chunkHeight,
            scale  : chunkScale,
            buffer : chunkBuffer,
        }
        currIndex += chunkWidth
    }

    return chunks
}

func newDrawChunk(chunk *ImageChunk, img image.Image, startX int, startY int, wg *sync.WaitGroup, conn net.Conn) {

    defer wg.Done()

    for x := chunk.xPos; x < (chunk.xPos + chunk.width); x++ {
        for y := 0; y < chunk.height; y++ {

            currPixelColor := img.At(x, y)
            if (currPixelColor != chunk.buffer.At(x, y)) {
                r, g, b, a := currPixelColor.RGBA()
                WritePixel(x + startX, y + startY, int(r>>8), int(g>>8), int(b>>8), int(a>>8), conn)
                chunk.buffer.Set(x + startX, y + startY, currPixelColor)
            }
        }
    }

}

func drawChunk(chunk *ImageChunk, img image.Image, startX int, startY int, wg *sync.WaitGroup, conn net.Conn) {

    defer wg.Done()

    for x := chunk.xPos; x < (chunk.xPos + chunk.width); x++ {
        for y := 0; y < chunk.height; y++ {
            scaledX := int(float64(x) / chunk.scale)
            scaledY := int(float64(y) / chunk.scale)

            // fmt.Println(x, y)
            // fmt.Println(chunk.xPos, chunk.width, chunk.height)
            currPixelColor := img.At(scaledX, scaledY)
            if (currPixelColor != chunk.buffer.At(scaledX, scaledY)) {
                r, g, b, a := currPixelColor.RGBA()
                WritePixel(x + startX, y + startY, int(r>>8), int(g>>8), int(b>>8), int(a>>8), conn)
                // fmt.Println("chunk.buffer BEFORE Set() = ", *chunk.buffer)
                chunk.buffer.Set(x + startX, y + startY, currPixelColor)
                // fmt.Println("chunk.buffer AFTER Set() = ", *chunk.buffer)
            } else {
                // fmt.Println("skipping pixel", currPixelColor)
            }
        }
    }

}
