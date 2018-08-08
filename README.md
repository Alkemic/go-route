# go-route - a regexp based router

``go-route`` allows you to create multi level regexp-based routing for your
application with registering named parameters in regexp.

The main change is that handler function definition changes from
``func(ResponseWriter, *Request)`` to ``func(ResponseWriter, *Request, map[string]string)``.
In given map are held params from URL.

## Middlewares

### Allowed methods

This middleware checks if request method is allowed.

Use consts from ``net/http`` for methos names, i.e. ``http.MethodPost``.

```go
package handlers

import (
    "net/http"

    "github.com/Alkemic/go-route/middleware"
)

middleware.AllowedMethods([]string{http.MethodPost, http.MethodPut}])(
    func(w ResponseWriter, r *Request, p map[string]string) {
        // actual code
    }
)
```

## Examples

### Example function

```go
package routing

import (
    ...
    "github.com/Alkemic/go-route"
)

func News(w ResponseWriter, r *Request, p map[string]string) {
    pk, _ := strconv.Atoi(p["pk"])
    news, _ := model.GetNews(pk)
    comments, _ := news.GetComments()
    _ = json.NewEncoder(w).Encode(map[string]interface{}{
        "Entry":    news,
        "Comments": comments,
    })
}

```

### Example routing definition

```go
package main

import (
    ...
    "github.com/Alkemic/go-route"
)

func main() {
    newsRoutes := route.RegexpRouter{}
    newsRoutes.Add(`^$`, view.NewsList)
    newsRoutes.Add(`^(?P<pk>\d+)/$`, view.News)
    newsRoutes.Add(`^(?P<pk>\d+),(?P<slug>[a-z\-_]+)\.html$`, view.News)

    routing := route.RegexpRouter{}
    routing.Add(`^/news/`, newsRoutes)

    // register custom 404 handler function
    route.NotFound = helper.Handle404

    log.Fatalln(http.ListenAndServe("0.0.0.0:80", routing.ServeHTTP))
}

```
