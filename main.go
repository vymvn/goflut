package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net"
	"os"
)


type canvasSize struct {
    width  int
    height int
}

func getSize(conn net.Conn) *canvasSize {

    var size canvasSize
    conn.Write([]byte("SIZE\n"))
    reply, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil {
        fmt.Println("Could not get size: ", err)
        return nil
    }

    fmt.Sscanf(reply, "SIZE %d %d", &size.width, &size.height)
    return &size
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

func drawImage(path string, startX int, startY int, conn net.Conn) error {
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
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b, a := img.At(x, y).RGBA()
			writePixel(x + startX, y + startY, int(r>>8), int(g>>8), int(b>>8), int(a>>8), conn)
		}
	}

	return nil
}

func main() {

    var host      *string = flag.String("host", "", "The PixelFlut server host ip or domain.")
    var port      *string = flag.String("port", "", "The port of the PixelFlut server.")
    var imagePath *string = flag.String("image", "", "The path to the image to draw.")

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
        fmt.Fprintln(os.Stderr, "Could not connect to \"" + connString + "\":\n\t", err)
        os.Exit(1)
    }
    defer conn.Close()

    drawImage(*imagePath, 100, 100, conn)

}
