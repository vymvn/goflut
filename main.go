package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)


func getSize(conn net.Conn) string {

    size := make([]byte, 1024)
    conn.Write([]byte("SIZE\n"))
    _, err := conn.Read(size);
    if err != nil {
        fmt.Println("Could not get size: ", err)
        return ""
    }
    return string(size)
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


func main() {

    var host *string = flag.String("host", "", "The PixelFlut server host ip or domain.")
    var port *string = flag.String("port", "", "The port of the PixelFlut server.")

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
        fmt.Println("Could not connect to \"" + connString + "\":\n\t", err)
        os.Exit(1)
    }
    defer conn.Close()

}
