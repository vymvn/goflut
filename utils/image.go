package utils

import (
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"net"
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


func getImageSize(img image.Image, conn net.Conn) (int, int) {

    bounds := img.Bounds()
    imgWidth := bounds.Max.X
    imgHeight := bounds.Max.Y

    return imgWidth, imgHeight
}


func ExpDrawImageThreaded(img image.Image, startX int, startY int, threads int, center bool, conn net.Conn) error {

    scaledWidth, scaledHeight := getImageSize(img, conn)

    if (center == true) {
        startX = (canvasSize.width / 2) - (scaledWidth / 2)
        startY = (canvasSize.height / 2) - (scaledHeight / 2)
    }

    chunkWidth := scaledWidth / threads
    // var chunks []*chunk = makeChunks(threads, chunkWidth, scaledHeight, scale) 
    var chunks []*expChunk = expMakeChunks(img, threads, chunkWidth, scaledHeight)

    bounds := image.Rect(0, 0, chunkWidth, scaledHeight)

    var wg sync.WaitGroup
    for i := 0; i < threads; i++ {
        wg.Add(1)

        draw.Draw(chunks[i].currBuffer, bounds, img, image.Point{chunks[i].xPos, 0}, draw.Src)
        go expDrawChunk(chunks[i], startX + chunks[i].xPos, startY, &wg, conn)
    }

    wg.Wait()

    return nil
}

func DrawImageThreaded(img image.Image, startX int, startY int, size float64, threads int, center bool, conn net.Conn) error {

    scaledWidth, scaledHeight, scale := getScaledImageSize(img, size, conn)

    if (center == true) {
        startX = (canvasSize.width / 2) - (scaledWidth / 2)
        startY = (canvasSize.height / 2) - (scaledHeight / 2)
    }

    chunkWidth := scaledWidth / threads
    var chunks []*chunk = makeChunks(threads, chunkWidth, scaledHeight, scale) 

    var wg sync.WaitGroup
    for i := 0; i < threads; i++ {
        wg.Add(1)

        go newDrawChunk(chunks[i], img, startX + chunks[i].xPos, startY, &wg, conn)
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

// func drawImageFromPath(path string, startX, startY int, threads int, size float64, conn net.Conn) error {
//
//     // t0 := time.Now()
//     f, err := os.Open(path)
//     if err != nil {
//         return err
//     }
//     defer f.Close()
//
//     img, _, err := image.Decode(f)
//     if err != nil {
//         return err
//     }
//
//     if (startX == -1 && startY == -1) {
//
//         scaledWidth, scaledHeight, _ := getImageSize(img, size, conn)
//
//         startX = (canvasSize.width / 2) - (scaledWidth / 2)
//         startY = (canvasSize.height / 2) - (scaledHeight / 2)
//     }
//
//     drawImage(img, startX, startY, threads, size, conn)
//     // fmt.Printf("drawImage runtime: %v\n", time.Since(t0))
//
//     return nil
// }

// func DrawImageFromPath(path string, startX, startY int, size float64, center bool, conn net.Conn) error {
//
//     f, err := os.Open(path)
//     if err != nil {
//         return err
//     }
//     defer f.Close()
//
//     img, _, err := image.Decode(f)
//     if err != nil {
//         return err
//     }
//
//     if (center == true) {
//
//         scaledWidth, scaledHeight, _ := getScaledImageSize(img, size, conn)
//         startX = (canvasSize.width / 2) - (scaledWidth / 2)
//         startY = (canvasSize.height / 2) - (scaledHeight / 2)
//     }
//
//     drawImage(img, startX, startY, size, center, conn)
//
//     return nil
// }
