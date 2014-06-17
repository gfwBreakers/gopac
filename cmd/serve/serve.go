package serve

import (
	"fmt"
	"net/http"
	"os"

	"github.com/codegangsta/cli"
)

type Static struct {
	file string
}

func (s *Static) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "HEAD" {
		return
	}
	f, _ := os.Open(s.file)
	defer f.Close()
	fi, _ := f.Stat()
	http.ServeContent(rw, r, s.file, fi.ModTime(), f)
}

func Action(c *cli.Context) {
	var port = c.String("port")
	s := &Static{"go.pac"}
	http.ListenAndServe(":"+port, s)
	fmt.Printf("Server's port is %s.\n", port)
}
