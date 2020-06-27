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

func subView(rw http.ResponseWriter, req *http.Request) {
	p := route.GetParams(req)
	fmt.Fprintf(rw, "sub params: %s", p)
}

func getHandle404(name string) func(rw http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, name+": Four OH! Four")
	}
}

func main() {
	flag.Parse()
	subRoutes := route.New()
	subRoutes.Add(`^/(?P<digit>\d+)$`, subView)
	subRoutes.Add(`^/(?P<param>.*)$`, subView)
	//subRoutes.NotFound = getHandle404("sub router")

	routes := route.New()
	routes.Add(`^/sub`, subRoutes)
	routes.Add(`^/$`, index)
	routes.Add(`^/(?P<param>[a-z]*)/(?P<param2>[a-z]*)/$`, view)
	routes.Add(`^/(?P<param>.*)/$`, view)
	routes.NotFound = getHandle404("main")

	log.Fatalln(http.ListenAndServe(*bindAddr, routes))
}
