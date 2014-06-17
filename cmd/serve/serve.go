package serve

import (
	"fmt"
	"net/http"
	"os"

	"github.com/codegangsta/cli"
)

type Tiny struct {
	file string
}

func (t *Tiny) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "HEAD" {
		return
	}
	f, _ := os.Open(t.file)
	defer f.Close()
	fi, _ := f.Stat()
	http.ServeContent(rw, r, t.file, fi.ModTime(), f)
}

func Action(c *cli.Context) {
	var port = c.String("port")
	if port == "" {
		return
	}

	t := Tiny{"go.pac"}
	http.ListenAndServe(":"+port, t)
}
