package utils

import (
	"image"
	"net"
)


func BouncingImage(img image.Image, x, y, xvel, yvel int, size float64, conn net.Conn) error {

    imgWidth, imgHeight, _ := getScaledImageSize(img, size, conn)

    drawCounter := 0
    for true {
        x += xvel
        y += yvel

        if x + imgWidth > canvasSize.width || x + imgWidth < 0 || x > canvasSize.width || x < 0 {
            xvel *= -1
        }
        if y + imgHeight > canvasSize.height || y + imgHeight < 0 || y > canvasSize.height || y < 0{
            yvel *= -1
        }

        if (drawCounter == 10) {
            DrawImage(img, x, y, size, false, conn)
            drawCounter = 0
        }
        drawCounter++
    }
    return nil
}
