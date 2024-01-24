package utils

import (
	"fmt"
	"image"
	"image/draw"
	"net"
	"os"

	"github.com/golang/freetype"
	"golang.org/x/image/font"
)


func DrawText(globalOpts *GlobalOptions, textOpts *TextOptions) error {

    connString := fmt.Sprintf("%s:%d", globalOpts.Host, globalOpts.Port)
    conn, err := net.Dial("tcp", connString)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Could not connect to \"" + connString + "\":\n", err)
        os.Exit(1)
    }
    defer conn.Close()

    getCanvasSize(&canvasSize, conn)

    fontBytes, err := os.ReadFile(textOpts.FontPath)
    if err != nil {
        // return err
        return fmt.Errorf("Couldn't open default font file\nUse -f to pass a .ttf font file: %w", err)
    }
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return err
	}

	// Initialize the context.
	// fg, bg := image.Black, image.White
	fg, bg := image.White, image.Black

    // fontFace := truetype.NewFace(f, &truetype.Options{Size: size})
	rgba := image.NewRGBA(image.Rect(0, 0, ((int(textOpts.FontSize) * len(textOpts.Text)) * 2), (int(textOpts.FontSize) * 6)))
	draw.Draw(rgba, rgba.Bounds(), bg, image.Point{}, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(300)
	c.SetFont(f)
	c.SetClip(image.Rect(0, 0, canvasSize.width, canvasSize.height))
	c.SetDst(rgba)
	c.SetSrc(fg)
	c.SetFontSize(textOpts.FontSize)
    // c.SetHinting(font.HintingNone)
    c.SetHinting(font.HintingFull)

    bounds := rgba.Bounds()
    width  := bounds.Max.X
    height := bounds.Max.Y

    pt := freetype.Pt(0, 0 +int(c.PointToFixed(textOpts.FontSize) >> 6))
    if _, err := c.DrawString(textOpts.Text, pt); err != nil {
        return err
    }


    // Use this after removing the scaling code or somthing
    // err = DrawImage(rgba, startX, startY, 1, center, conn)
    // if err != nil {
    //     return err
    // }

    if (textOpts.Center == true) {
        globalOpts.StartX = (canvasSize.width  / 2) - (width  / 2)
        globalOpts.StartY = (canvasSize.height / 2) - (height / 2)
    }

    for x := 0; x < width; x++ {
        for y := 0; y < height; y++ {
            r, g, b, a := rgba.At(x, y).RGBA()
            WritePixel(x + globalOpts.StartX, y + globalOpts.StartY, int(r>>8), int(g>>8), int(b>>8), int(a>>8), conn)
        }
    }

    return nil

}
