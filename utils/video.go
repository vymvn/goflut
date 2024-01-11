package utils

import (
	"image"
	"net"

	"github.com/AlexEidt/Vidio"
)


func DrawVideo(videoPath string, startX, startY int, size float64, center bool, conn net.Conn) {
    video, _ := vidio.NewVideo(videoPath)

    img := image.NewRGBA(image.Rect(0, 0, video.Width(), video.Height()))
    video.SetFrameBuffer(img.Pix)

    frame := 0
    for video.Read() {
        DrawImage(img, startX, startY, size, center, conn)
        frame++
    }
}
