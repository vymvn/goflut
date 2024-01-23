package utils

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"net"
	"os"
	"sync"
)

func getScaledImageSize(img image.Image, size float64, conn net.Conn) (int, int, float64) {

    bounds := img.Bounds()
    imgWidth := bounds.Max.X
    imgHeight := bounds.Max.Y

    getCanvasSize(&canvasSize, conn)

    // Calculate scaling factors
    widthScale := float64(canvasSize.width) / float64(imgWidth)
    heightScale := float64(canvasSize.height) / float64(imgHeight)
    scale := math.Min(widthScale, heightScale) * size

    // Calculate scaled dimensions
    scaledWidth := int(float64(imgWidth) * scale)
    scaledHeight := int(float64(imgHeight) * scale)

    return scaledWidth, scaledHeight, scale
}

func getImageSize(img image.Image) (int, int) {

    bounds := img.Bounds()
    imgWidth := bounds.Max.X
    imgHeight := bounds.Max.Y

    return imgWidth, imgHeight
}

func ExpDrawImageThreaded(chunks []*ImageChunk, globalOpts *GlobalOptions) error {

    var wg sync.WaitGroup
    for i := 0; i < len(chunks); i++ {
        wg.Add(1)

        // draw.Draw(chunks[i].currBuffer, bounds, img, image.Point{chunks[i].xPos, 0}, draw.Src)
        // go expDrawChunk(chunks[i], startX + chunks[i].xPos, startY, &wg, newConn)
        // go drawChunk(chunks[i], img, startX, startY, &wg, chunks[i].conn)
        go expDrawImageChunk(chunks[i], globalOpts.StartX, globalOpts.StartY, &wg)
    }

    wg.Wait()

    // for i := 0; i < len(chunks); i++ {
    //     chunks[i].conn.Close()
    // }
    
    return nil
}

func DrawImageThreaded(img image.Image, startX int, startY int, size float64, threads int, center bool, conn net.Conn) error {

    scaledWidth, scaledHeight, scale := getScaledImageSize(img, size, conn)

    conn.Close()
    if (center == true) {
        startX = (canvasSize.width / 2) - (scaledWidth / 2)
        startY = (canvasSize.height / 2) - (scaledHeight / 2)
    }

    chunkWidth := scaledWidth / threads
    var chunks []*ImageChunk = makeChunks(threads, chunkWidth, scaledHeight, scale) 

    var wg sync.WaitGroup
    for i := 0; i < threads; i++ {
        wg.Add(1)

        connString := fmt.Sprintf("localhost:1234")
        // connString := fmt.Sprintf("pixelflut.uwu.industries:1234")
        newConn, err := net.Dial("tcp", connString)
        if err != nil {
            fmt.Fprintln(os.Stderr, "Could not connect to \"" + connString + "\":\n", err)
            os.Exit(1)
        }
        defer newConn.Close()

        // connString := fmt.Sprintf("localhost:1234")
        // draw.Draw(chunks[i].currBuffer, bounds, img, image.Point{chunks[i].xPos, 0}, draw.Src)
        // go expDrawChunk(chunks[i], startX + chunks[i].xPos, startY, &wg, newConn)
        go drawChunk(chunks[i], img, startX, startY, &wg, newConn)
        // go newDrawChunk(chunks[i], img, startX + chunks[i].xPos, startY, &wg, conn)
    }

    wg.Wait()

    return nil
}

func DrawImage(img image.Image, startX int, startY int, size float64, center bool, conn net.Conn) error {

    scaledWidth, scaledHeight, scale := getScaledImageSize(img, size, conn)

    if (center == true) {
        startX = (canvasSize.width / 2) - (scaledWidth / 2)
        startY = (canvasSize.height / 2) - (scaledHeight / 2)
    }

    for x := 0; x < scaledWidth; x++ {

        for y := 0; y <= scaledHeight; y++ {
            scaledX := int(float64(x) / scale)
            scaledY := int(float64(y) / scale)

            r, g, b, a := img.At(scaledX, scaledY).RGBA()
            err := WritePixel(x + startX, y + startY, int(r>>8), int(g>>8), int(b>>8), int(a>>8), conn)
            if (err != nil) {
                return err
            }
        }
    }


    return nil
}
