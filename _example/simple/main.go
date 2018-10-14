package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/Alkemic/go-route"
)

var (
	bindAddr = flag.String("bind", "localhost:8080", "address to bind to")
)

func index(rw http.ResponseWriter, req *http.Request, p map[string]string) {
	fmt.Fprint(rw, "index")
}

func view(rw http.ResponseWriter, req *http.Request, p map[string]string) {
	fmt.Fprintf(rw, "param: %s", p["param"])
}

func main() {
	flag.Parse()
	routes := route.RegexpRouter{}
	routes.Add(`^/$`, index)
	routes.Add(`^/(?P<param>.*)$`, view)
	log.Fatalln(http.ListenAndServe(*bindAddr, routes))
}
