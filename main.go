package main

import (
	"os"
	"runtime"

	"github.com/codegangsta/cli"
	"github.com/gfwBreakers/gopac/cmd/build"
	"github.com/gfwBreakers/gopac/cmd/serve"
)

const APP_VER = "0.0.0"

var app = cli.NewApp()

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	app.Name = "gopac"
	app.Usage = "Generate proxy auto-config rules and host them."
	app.Version = APP_VER
	app.Commands = append(app.Commands,
		cli.Command{
			Name:   "build",
			Usage:  "Generate proxy auto-config rules",
			Action: build.Action,
			Flags: []cli.Flag{
				cli.StringFlag{"proxy, x", "SOCKS5 127.0.0.1:8964; SOCKS 127.0.0.1:8964; DIRECT", "Examples: SOCKS5 127.0.0.1:8964; SOCKS 127.0.0.1:8964; PROXY 127.0.0.1:6489"},
			},
		},
		cli.Command{
			Name:   "serve",
			Usage:  "Start pac server",
			Action: serve.Action,
			Flags: []cli.Flag{
				cli.StringFlag{"port, p", "0", "Pac Server Port [OPTIONAL], examples: 8970"},
			},
		},
	)
}

func main() {
	app.Run(os.Args)
}
