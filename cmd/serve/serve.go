package serve

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"

	"github.com/codegangsta/cli"
)

func server(port string) {
	fmt.Printf("Port is %s.\n", port)
	ln, _ := net.Listen("tcp", ":"+port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	f, _ := os.Open("go.pac")
	defer f.Close()
	io.Copy(conn, f)
}

func Action(c *cli.Context) {
	var port = c.Int("port")
	if port <= 0 {
		return
	}

	server(strconv.Itoa(port))
}
