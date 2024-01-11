package utils

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"net"
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

// func drawImage(img image.Image, startX int, startY int, threads int, size float64, conn net.Conn) error {
//
//     scaledWidth, scaledHeight, scale := getImageSize(img, size, conn)
//
//     chunkWidth := int(scaledWidth / threads)
//     var chunks []chunk = makeChunks(threads, chunkWidth, scaledHeight, scale) 
//
//     var wg sync.WaitGroup
//     for i := 0; i < threads; i++ {
//         wg.Add(1)
//         go drawChunk(chunks[i], img, startX, startY, &wg, conn)
//     }
//
//     wg.Wait()
//
//     return nil
// }

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
