package utils

import (
	"image"
	"image/color"
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

func NewDrawImage(img [][]color.RGBA, numThreads int) {
    rowsPerThread := len(img) / numThreads

    var wg sync.WaitGroup

    for i := 0; i < numThreads; i++ {
        startRow := i * rowsPerThread
        endRow := startRow + rowsPerThread

        // The last thread may have extra rows if the image height is not divisible evenly
        if i == numThreads-1 {
            endRow = len(img)
        }

        wg.Add(1)
        go renderSection(img, startRow, endRow, &wg)
    }

    // Wait for all goroutines to finish
    wg.Wait()
}

func renderSection(img [][]color.RGBA, startRow, endRow int, wg *sync.WaitGroup) {
	defer wg.Done()

	// Render the specified section of the image
	for row := startRow; row < endRow; row++ {
		for col := 0; col < len(img[row]); col++ {
			// Your rendering logic here
			img[row][col] = color.RGBA{255, 0, 0, 255} // Example: Set pixel color to red
		}
	}
}

func DrawImageThreaded(img image.Image, startX int, startY int, size float64, threads int, center bool, conn net.Conn) error {

    scaledWidth, scaledHeight, scale := getScaledImageSize(img, size, conn)

    chunkWidth := scaledWidth / threads
    var chunks []*chunk = makeChunks(threads, chunkWidth, scaledHeight, scale) 

    var wg sync.WaitGroup
    for i := 0; i < threads; i++ {
        wg.Add(1)
        go drawChunk(chunks[i], img, startX, startY, &wg, conn)
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
