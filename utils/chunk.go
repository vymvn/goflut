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
    buffer image.RGBA
}


func makeChunks(threadsCount int, chunkWidth int, chunkHeight int, chunkScale float64) []chunk {

    chunks := make([]chunk, threadsCount)   // As many chunks as threads

    currIndex := 0
    for i := 0; i < len(chunks); i++{

        chunkBuffer := image.NewRGBA(image.Rect(0, 0, chunkWidth, chunkHeight))
        chunks[i] = chunk{
            xPos   : currIndex,
            width  : chunkWidth,
            height : chunkHeight,
            scale  : chunkScale,
            buffer : *chunkBuffer,
        }

        currIndex += chunkWidth
    }

    return chunks
}

func drawChunk(chunk *chunk, img image.Image, startX int, startY int, wg *sync.WaitGroup, conn net.Conn) {

    defer wg.Done()

    for x := chunk.xPos; x < (chunk.xPos + chunk.width); x++ {
        for y := 0; y < chunk.height; y++ {
            scaledX := int(float64(x) / chunk.scale)
            scaledY := int(float64(y) / chunk.scale)

            currPixelColor := img.At(scaledX, scaledY)
            if (currPixelColor != chunk.buffer.At(scaledX, scaledY)) {
                r, g, b, a := currPixelColor.RGBA()
                WritePixel(x + startX, y + startY, int(r>>8), int(g>>8), int(b>>8), int(a>>8), conn)
                chunk.buffer.Set(x, x, currPixelColor)
            }
        }
    }

}
