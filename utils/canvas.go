package utils

import (
	"bufio"
	"fmt"
	"net"
)

type Size struct {
    width  int
    height int
}

var (
    canvasSize Size
)

func getCanvasSize(canvasSize *Size, conn net.Conn) error {

    conn.Write([]byte("SIZE\n"))
    reply, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil {
        return err
    }

    fmt.Sscanf(reply, "SIZE %d %d", &canvasSize.width, &canvasSize.height)

    return nil
}
