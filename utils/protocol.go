package utils

import (
	"fmt"
	"net"
)


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
