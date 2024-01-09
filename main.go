package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	// "image/color"
	// "image/draw"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"net"
	"os"
	"sync"

	"github.com/golang/freetype"
	"golang.org/x/image/font"
	// "time"
)


type size struct {
    width  int
    height int
}

type chunk struct {
    xPos   int
    width  int
    height int
    scale  float64
}


var (

    host      *string = flag.String("host", "", "The PixelFlut server host ip or domain.")
    port      *string = flag.String("port", "", "The port of the PixelFlut server.")
    imagePath *string = flag.String("image", "", "The path to the image to draw.")
    threads      *int = flag.Int("threads", 1, "Number of threads to use.")

    canvasSize size
)

func getCanvasSize(conn net.Conn) error {

    conn.Write([]byte("SIZE\n"))
    reply, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil {
        return err
    }

    fmt.Sscanf(reply, "SIZE %d %d", &canvasSize.width, &canvasSize.height)

    return nil
}

func writePixel(x, y, r, g, b, a int, conn net.Conn) {
    var cmd string
    if a == 255 {
        cmd = fmt.Sprintf("PX %d %d %02x%02x%02x\n", x, y, r, g, b)
    } else {
        cmd = fmt.Sprintf("PX %d %d %02x%02x%02x%02x\n", x, y, r, g, b, a)
    }
    conn.Write([]byte(cmd))
}

func drawRect(x, y, w, h, r, g, b int, conn net.Conn) {
    for i := x; i < x+w; i++ {
        for j := y; j < y+h; j++ {
            writePixel(i, j, r, g, b, 255, conn)
        }
    }
}

func bouncingImage(x, y, xvel, yvel int, path string, size float64, conn net.Conn, threads int) error {

    f, err := os.Open(path)
    if err != nil {
        return err
    }
    defer f.Close()

    img, _, err := image.Decode(f)
    if err != nil {
        return err
    }

    imgWidth, imgHeight, _ := getImageSize(img, size, conn)

    for true {
        x += xvel
        y += yvel

        if x + imgWidth > canvasSize.width || x + imgWidth < 0 || x > canvasSize.width || x < 0 {
            xvel *= -1
        }
        if y + imgHeight > canvasSize.height || y + imgHeight < 0 || y > canvasSize.height || y < 0{
            yvel *= -1
        }

        drawImage(path, x, y, threads, size, conn)
    }
    return nil
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
            writePixel(x + startX, y + startY, int(r>>8), int(g>>8), int(b>>8), int(a>>8), conn)
        }
    }

}

func getImageSize(img image.Image, size float64, conn net.Conn) (int, int, float64) {

    bounds := img.Bounds()
    imgWidth := bounds.Max.X
    imgHeight := bounds.Max.Y

    // Calculate scaling factors
    widthScale := float64(canvasSize.width) / float64(imgWidth)
    heightScale := float64(canvasSize.height) / float64(imgHeight)
    scale := math.Min(widthScale, heightScale) * size

    // Calculate scaled dimensions
    scaledWidth := int(float64(imgWidth) * scale)
    scaledHeight := int(float64(imgHeight) * scale)

    return scaledWidth, scaledHeight, scale
}

func processImage(img image.Image, startX int, startY int, threads int, size float64, conn net.Conn) error {

    scaledWidth, scaledHeight, scale := getImageSize(img, size, conn)

    chunkWidth := int(scaledWidth / threads)
    var chunks []chunk = makeChunks(threads, chunkWidth, scaledHeight, scale) 

    var wg sync.WaitGroup
    for i := 0; i < threads; i++ {
        wg.Add(1)
        go drawChunk(chunks[i], img, startX, startY, &wg, conn)
    }

    wg.Wait()

    return nil
}

func drawImage(path string, startX, startY int, threads int, size float64, conn net.Conn) error {

    // t0 := time.Now()
    f, err := os.Open(path)
    if err != nil {
        return err
    }
    defer f.Close()

    img, _, err := image.Decode(f)
    if err != nil {
        return err
    }

    if (startX == -1 && startY == -1) {

        scaledWidth, scaledHeight, _ := getImageSize(img, size, conn)

        startX = (canvasSize.width / 2) - (scaledWidth / 2)
        startY = (canvasSize.height / 2) - (scaledHeight / 2)
    }

    processImage(img, startX, startY, threads, size, conn)
    // fmt.Printf("drawImage runtime: %v\n", time.Since(t0))

    return nil
}

func drawText(text string, startX, startY int, size float64, conn net.Conn) {

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
	rgba := image.NewRGBA(image.Rect(0, 0, 800, 200))
	// draw.Draw(rgba, rgba.Bounds(), bg, image.Point{}, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(300)
	c.SetFont(f)
	c.SetFontSize(size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
    c.SetHinting(font.HintingNone)
    // c.SetHinting(font.HintingFull)

    pt := freetype.Pt(10, 5 +int(c.PointToFixed(size) >> 6))

    if _, err := c.DrawString(text, pt); err != nil {
        log.Println(err)
        return
    }

    bounds := rgba.Bounds()
    width  := bounds.Max.X
    height := bounds.Max.Y

    // processImage(rgba, -1, -1, 1, 1, conn)
    for x := 0; x < width; x++ {
        for y := 0; y < height; y++ {
            r, g, b, a := rgba.At(x, y).RGBA()
            writePixel(x + startX, y + startY, int(r>>8), int(g>>8), int(b>>8), int(a>>8), conn)
        }
    }

}

func wipeCanvas(conn net.Conn) {
    drawRect(0, 0, canvasSize.width, canvasSize.height, 50, 50, 50, conn)
}

func applyBackground(r, g, b int, conn net.Conn) {
    drawRect(0, 0, canvasSize.width, canvasSize.height, r, g, b, conn)
}

func main() {

    required := []string{"host", "port"}
    flag.Parse()

    seen := make(map[string]bool)
    flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })
    for _, req := range required {
        if !seen[req] {
            flag.Usage()
            os.Exit(2)
        }
    }

    if (*threads == 0) {
        fmt.Fprintln(os.Stderr, "Don't be silly")
        os.Exit(1)
    }

    connString := fmt.Sprintf("%s:%s", *host, *port)
    conn, err := net.Dial("tcp", connString)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Could not connect to \"" + connString + "\":\n", err)
        os.Exit(1)
    }
    defer conn.Close()

    getCanvasSize(conn)
    wipeCanvas(conn)
    drawText("text rendering from the hood", 100, 100, 12, conn)

    // bouncingImage(0, 0, 50, 70, *imagePath, 0.2, conn, *threads)

    // err = drawImage(*imagePath, -1, -1, *threads, 0.5, conn)
    // if err != nil {
    //     fmt.Fprintln(os.Stderr, "Could not draw image:" + "\n", err)
    //     os.Exit(1)
    // }
}
