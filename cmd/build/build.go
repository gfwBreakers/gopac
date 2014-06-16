package build

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/codegangsta/cli"
	"github.com/go-libs/iputils"
)

var (
	F      uint32 = 0xffffffff
	EOL           = byte('\n')
	URL           = "http://ftp.apnic.net/apnic/stats/apnic/delegated-apnic-latest"
	CNIPV4        = regexp.MustCompile(`apnic\|(CN|cn)\|ipv4\|([0-9\.]+)\|([0-9]+)\|([0-9]+)\|a.*`)
)

type Node struct {
	IP, Mask, Mask2 uint32
}

func NewNode(ip, mask string, mask2 uint32) Node {
	return Node{iputils.IP2Long(ip), iputils.IP2Long(mask), mask2}
}

type Graph []Node

func (g Graph) Len() int           { return len(g) }
func (g Graph) Swap(i, j int)      { g[i], g[j] = g[j], g[i] }
func (g Graph) Less(i, j int) bool { return g[i].IP < g[j].IP }

func fetchIPData(results *Graph) (err error) {
	var (
		n   int
		buf string
		r   *bufio.Reader

		startIP, prevIP, smask string
		numIP                  int
		imask                  uint32
	)

	fmt.Println("Fetching data from apnic.net, it might take a few minutes, please wait...")

	res, err := http.Get(URL)
	if err != nil {
		panic(err)
	}
	r = bufio.NewReader(res.Body)
	defer res.Body.Close()
	for {
		// read line by line
		buf, err = r.ReadString(EOL)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		matches := CNIPV4.FindStringSubmatch(buf)
		if len(matches) > 0 {
			n++
			startIP = matches[2]
			numIP, _ = strconv.Atoi(matches[3])

			imask = F ^ uint32(numIP-1)
			smask = fmt.Sprintf("%02x", imask)
			mask := [4]string{}
			mask[0] = smask[0:2]
			mask[1] = smask[2:4]
			mask[2] = "0"
			mask[3] = "0"

			for i, s := range mask[:2] {
				num, _ := strconv.ParseInt(s, 16, 10)
				mask[i] = fmt.Sprintf("%d", num)
			}

			mask2 := 32 - uint32(math.Log2(float64(numIP)))

			ip := strings.Split(startIP, ".")
			ip[2] = "0"
			ip[3] = "0"
			startIP = strings.Join(ip, ".")

			maskIP := fmt.Sprintf("%s.%s.%s.%s", mask[0], mask[1], mask[2], mask[3])

			if startIP != prevIP {
				*results = append(*results, NewNode(startIP, maskIP, mask2))
				prevIP = startIP
			}
		}
	}
	return nil
}

func Action(c *cli.Context) {
	var pacfile = "go.pac"
	var results = make(Graph, 0)
	results = append(results, NewNode("127.0.0.1", "255.0.0.0", 0))
	results = append(results, NewNode("10.0.0.0", "255.0.0.0", 0))
	results = append(results, NewNode("127.0.0.1", "255.240.0.0", 0))
	results = append(results, NewNode("192.168.0.0", "255.255.0.0", 0))
	fetchIPData(&results)
	sort.Sort(results)

	file, err := os.Create(pacfile)
	if err != nil {
		panic(err)
	}

	t, err := template.ParseFiles("templates/pac.tmpl")
	data := make(map[string]interface{}, 0)
	data["Graph"] = results
	data["Proxy"] = c.String("proxy")
	err = t.Execute(file, data)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Rules: %d items.\n", len(results))
	fmt.Printf("Usage: Use the newly created %s as your web browser's automatic \n", pacfile)
	fmt.Printf("PAC(Proxy auto-config) file.\n")
}
