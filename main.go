package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"net"
	"os"
	"sync"
)


type size struct {
    width  int
    height int
}

func getSize(conn net.Conn) *size {

    var canvasSize size
    conn.Write([]byte("SIZE\n"))
    reply, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil {
        fmt.Println("Could not get size: ", err)
        return nil
    }

    fmt.Sscanf(reply, "SIZE %d %d", &canvasSize.width, &canvasSize.height)
    return &canvasSize
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

func drawImage(path string, startX int, startY int, threads int, conn net.Conn) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return err
	}

	bounds := img.Bounds()
	imgWidth := bounds.Max.X
	imgHeight := bounds.Max.Y

    canvasSize := getSize(conn)

	// Calculate scaling factors
	widthScale := float64(canvasSize.width) / float64(imgWidth)
	heightScale := float64(canvasSize.height) / float64(imgHeight)
	scale := math.Min(widthScale, heightScale)

	// Calculate scaled dimensions
	scaledWidth := int(float64(imgWidth) * scale)
	scaledHeight := int(float64(imgHeight) * scale)

	var wg sync.WaitGroup
	work := make(chan int, scaledWidth)
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for x := range work {
				for y := 0; y < scaledHeight; y++ {
					scaledX := int(float64(x) / scale)
					scaledY := int(float64(y) / scale)

					r, g, b, a := img.At(scaledX, scaledY).RGBA()
					writePixel(x + startX, y + startY, int(r>>8), int(g>>8), int(b>>8), int(a>>8), conn)
				}
			}
		}()
	}

	for x := 0; x < scaledWidth; x++ {
		work <- x
	}
	close(work)

	wg.Wait()

	return nil
}

func main() {

    var host      *string = flag.String("host", "", "The PixelFlut server host ip or domain.")
    var port      *string = flag.String("port", "", "The port of the PixelFlut server.")
    var imagePath *string = flag.String("image", "", "The path to the image to draw.")
    var threads      *int = flag.Int("threads", 4, "Number of threads to use.")

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

    connString := fmt.Sprintf("%s:%s", *host, *port)
    conn, err := net.Dial("tcp", connString)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Could not connect to \"" + connString + "\":\n", err)
        os.Exit(1)
    }
    defer conn.Close()

    err = drawImage(*imagePath, 0, 0, *threads, conn)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Could not draw image:" + "\n", err)
        os.Exit(1)
    }

}
