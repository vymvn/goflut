package utils

import (
	"fmt"
	"math/rand"
	"net"
)

func makeConnection(connString string) (error, net.Conn) {

    newConn, err := net.Dial("tcp", connString)
    if err != nil {
        return err, nil
    }

    return nil, newConn
}

func WritePixel(x, y, r, g, b, a int, conn net.Conn) error{
    var cmd string
    if a == 255 {
        cmd = fmt.Sprintf("PX %d %d %02x%02x%02x\n", x, y, r, g, b)
    } else {
        cmd = fmt.Sprintf("PX %d %d %02x%02x%02x%02x\n", x, y, r, g, b, a)
    }
    _, err := conn.Write([]byte(cmd))
    if (err != nil) {
        return err
    }

    return nil
}

func drawRect(x, y, w, h, r, g, b, a int, conn net.Conn) {
    for i := x; i < x+w; i++ {
        for j := y; j < y+h; j++ {
            WritePixel(i, j, r, g, b, a, conn)
        }
    }
}

func Noise(startX, startY int, connString string) error {

    err, conn := makeConnection(connString)
    if err != nil {
        return err
    }
    defer conn.Close()
    getCanvasSize(&canvasSize, conn)

    for x := 0; x < canvasSize.width; x++ {

        for y := 0; y < canvasSize.height; y++ {
            WritePixel(x + startX, y + startY, rand.Intn(256), rand.Intn(256), rand.Intn(256), 255, conn)
        }
    }

    return nil
}
