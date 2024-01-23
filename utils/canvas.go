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

func WipeCanvas(connString string) error {
    err, conn := makeConnection(connString)
    if err != nil {
        return err
    }
    defer conn.Close()

    getCanvasSize(&canvasSize, conn)
    drawRect(0, 0, canvasSize.width, canvasSize.height, 50, 50, 50, 255, conn)

    return nil
}

func ApplyBackground(r, g, b int, connString string) error {
    err, conn := makeConnection(connString)
    if err != nil {
        return err
    }
    defer conn.Close()

    getCanvasSize(&canvasSize, conn)
    drawRect(0, 0, canvasSize.width, canvasSize.height, r, g, b, 255, conn)

    return nil
}
