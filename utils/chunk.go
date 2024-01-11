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
}


func makeChunks(threadsCount int, chunkWidth int, chunkHeight int, chunkScale float64) []chunk {

    chunks := make([]chunk, threadsCount)   // As many chunks as threads
    currIndex := 0
    for i := 0; i < len(chunks); i++{
        chunks[i] = chunk{
            xPos   : currIndex,
            width  : chunkWidth,
            height : chunkHeight,
            scale  : chunkScale,
        }

        currIndex += chunkWidth
    }

    return chunks
}

func drawChunk(chunk chunk, img image.Image, startX int, startY int, wg *sync.WaitGroup, conn net.Conn) {

    defer wg.Done()

    for x := chunk.xPos; x < (chunk.xPos + chunk.width); x++ {
        for y := 0; y <= chunk.height; y++ {
            scaledX := int(float64(x) / chunk.scale)
            scaledY := int(float64(y) / chunk.scale)

            r, g, b, a := img.At(scaledX, scaledY).RGBA()
            WritePixel(x + startX, y + startY, int(r>>8), int(g>>8), int(b>>8), int(a>>8), conn)
        }
    }

}
