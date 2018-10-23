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

func index(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprint(rw, "index")
}

func view(rw http.ResponseWriter, req *http.Request) {
	p := route.GetParams(req)
	fmt.Fprintf(rw, "params: %s", p)
}

func main() {
	flag.Parse()
	subRoutes := route.RegexpRouter{}
	subRoutes.Add(`/(?P<param>.*)$`, view)

	routes := route.RegexpRouter{}
	routes.Add(`^/sub/(?P<digit>\d+)`, subRoutes)
	routes.Add(`^/$`, index)
	routes.Add(`^/(?P<param>[a-z]*)/(?P<param2>[a-z]*)/$`, view)
	routes.Add(`^/(?P<param>.*)$`, view)

	log.Fatalln(http.ListenAndServe(*bindAddr, routes))
}
