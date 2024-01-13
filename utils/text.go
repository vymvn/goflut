package utils

import (
	"image"
	"image/draw"
	"net"
	"os"

	"github.com/golang/freetype"
	"golang.org/x/image/font"
)


func DrawText(text string, startX, startY int, size float64, center bool, conn net.Conn) error {

    getCanvasSize(&canvasSize, conn)

    fontBytes, err := os.ReadFile("fonts/Lato-Regular.ttf")
    if err != nil {
        return err
    }
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return err
	}

	// Initialize the context.
	fg, bg := image.Black, image.White

	rgba := image.NewRGBA(image.Rect(0, 0, ((int(size) * len(text)) * 2) + 10, (int(size) * 6)))
	draw.Draw(rgba, rgba.Bounds(), bg, image.Point{}, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(300)
	c.SetFont(f)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	c.SetFontSize(size)
    // c.SetHinting(font.HintingNone)
    c.SetHinting(font.HintingFull)

    bounds := rgba.Bounds()
    width  := bounds.Max.X
    height := bounds.Max.Y

    pt := freetype.Pt(5, 0 +int(c.PointToFixed(size) >> 6))
	c.SetSrc(image.Black)
    if _, err := c.DrawString(text, pt); err != nil {
        return err
    }


    // Use this after removing the scaling code or somthing
    // err = DrawImage(rgba, startX, startY, 1, center, conn)
    // if err != nil {
    //     return err
    // }

    if (center == true) {
        startX = (canvasSize.width  / 2) - (width  / 2)
        startY = (canvasSize.height / 2) - (height / 2)
    }

    for x := 0; x < width; x++ {
        for y := 0; y < height; y++ {
            r, g, b, a := rgba.At(x, y).RGBA()
            WritePixel(x + startX, y + startY, int(r>>8), int(g>>8), int(b>>8), int(a>>8), conn)
        }
    }

    return nil

}
