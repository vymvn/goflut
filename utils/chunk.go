package utils

import (
	"image"
	"net"
	"sync"
)


type chunk struct {
    xPos   int
    width  int
    height int
    scale  float64
    buffer *image.RGBA
}


// func makeChunks(threadsCount int, chunkWidth int, chunkHeight int, chunkScale float64) []chunk {
//
//     chunks := make([]chunk, threadsCount)   // As many chunks as threads
//
//     currIndex := 0
//     for i := 0; i < len(chunks); i++{
//
//         chunkBuffer := image.NewRGBA(image.Rect(0, 0, chunkWidth, chunkHeight))
//         chunks[i] = chunk{
//             xPos   : currIndex,
//             width  : chunkWidth,
//             height : chunkHeight,
//             scale  : chunkScale,
//             buffer : chunkBuffer,
//         }
//
//         currIndex += chunkWidth
//     }
//
//     return chunks
// }

func makeChunks(threadsCount int, chunkWidth int, chunkHeight int, chunkScale float64) []*chunk {

    chunks := make([]*chunk, threadsCount)   // As many chunks as threads

    currIndex := 0
    for i := 0; i < len(chunks); i++{

        chunkBuffer := image.NewRGBA(image.Rect(0, 0, chunkWidth, chunkHeight))
        chunks[i] = &chunk{
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


// func drawChunk(chunk *chunk, img image.Image, startX int, startY int, wg *sync.WaitGroup, conn net.Conn) {
//
//     defer wg.Done()
//
//     fmt.Println("Before: ", img.At(0, 0), chunk.buffer.At(0, 0))
//     for x := chunk.xPos; x < (chunk.xPos + chunk.width); x++ {
//         for y := 0; y < chunk.height; y++ {
//
//             // fmt.Println(x, y)
//             // fmt.Println(chunk.xPos, chunk.width, chunk.height)
//             currPixelColor := img.At(x, y)
//             // fmt.Println("chunk.buffer BEFORE if statement = ", chunk.buffer.At(x, y))
//             // fmt.Println(currPixelColor, chunk.buffer.At(x, y))
//             if (currPixelColor != chunk.buffer.At(x, y)) {
//                 r, g, b, a := currPixelColor.RGBA()
//                 WritePixel(x + startX, y + startY, int(r>>8), int(g>>8), int(b>>8), int(a>>8), conn)
//                 // fmt.Println("currPixelColor BEFORE Set() = ", currPixelColor)
//                 // fmt.Println("chunk.buffer BEFORE Set() = ", chunk.buffer.At(x, y))
//                 chunk.buffer.Set(x, y, currPixelColor)
//                 // fmt.Println("chunk.buffer AFTER Set() = ", chunk.buffer.At(x, y), "\n")
//             } else {
//                 fmt.Println("Skipping already drawn pixel: ", currPixelColor)
//             }
//         }
//     }
//
//     fmt.Println("After: ", img.At(0, 0), chunk.buffer.At(0, 0))
//
//
// }

func drawChunk(chunk *chunk, img image.Image, startX int, startY int, wg *sync.WaitGroup, conn net.Conn) {

    defer wg.Done()

    for x := chunk.xPos; x < (chunk.xPos + chunk.width); x++ {
        for y := 0; y <= chunk.height; y++ {
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
