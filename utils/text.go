package utils

import (
	"image"
	"log"
	"net"
	"os"

	"github.com/golang/freetype"
	"golang.org/x/image/font"
)


func DrawText(text string, startX, startY int, size float64, color string, center bool, conn net.Conn) {

    getCanvasSize(&canvasSize, conn)

    fontBytes, err := os.ReadFile("fonts/Lato-Regular.ttf")
    if err != nil {
        log.Println(err)
        return
    }
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}

	// Initialize the context.
	// fg, bg := image.Black, image.White
	fg := image.Black

    switch color {

    case "white":
	    fg = image.White

    case "black":
        fg = image.Black

    default:
        fg = image.Black
    }

	rgba := image.NewRGBA(image.Rect(0, 0, canvasSize.width, canvasSize.height))
	// draw.Draw(rgba, rgba.Bounds(), bg, image.Point{}, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(300)
	c.SetFont(f)
	c.SetFontSize(size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
    // c.SetHinting(font.HintingNone)
    c.SetHinting(font.HintingFull)

    pt := freetype.Pt(10, 5 +int(c.PointToFixed(size) >> 6))
	c.SetSrc(image.Black)
    if _, err := c.DrawString(text, pt); err != nil {
        log.Println(err)
        return
    }


    pt = freetype.Pt(8 , 4 +int(c.PointToFixed(size - 2) >> 6))
	c.SetSrc(image.White)
    if _, err := c.DrawString(text, pt); err != nil {
        log.Println(err)
        return
    }

    bounds := rgba.Bounds()
    width  := bounds.Max.X
    height := bounds.Max.Y

    if (center == true) {
        startX = (canvasSize.width / 2) - (width / 2)
        startY = (canvasSize.height / 2) - (height / 2)
    }

    for x := 0; x < width; x++ {
        for y := 0; y < height; y++ {
            r, g, b, a := rgba.At(x, y).RGBA()
            WritePixel(x + startX, y + startY, int(r>>8), int(g>>8), int(b>>8), int(a>>8), conn)
        }
    }

}
