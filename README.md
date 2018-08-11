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
package main

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

### Panic interceptors

This middleware catches panics and logs them, and then returns 500 to the user. 500 handler can be overriden by

Differences between ``PanicInterceptor`` and ``PanicInterceptorWithLogger`` is that, the second one is parametrised
decorator that accepts instance of ``*log.Logger``. If first one is use, it uses standard ``log.Printf``/``log.Println``.


```go
package main

import (
    "log"
    "net/http"
    "os"

    "github.com/Alkemic/go-route/middleware"
)

middleware.PanicInterceptor(
    func(w ResponseWriter, r *Request, p map[string]string) {
        // actual code
    }
)
```

```go
package main

import (
    "log"
    "net/http"
    "os"

    "github.com/Alkemic/go-route/middleware"
)

var logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

middleware.PanicInterceptorWithLogger(logger)(
    func(w ResponseWriter, r *Request, p map[string]string) {
        // actual code
    }
)
```

### Set headers

This middleware is used to set given headers to response.

```go
package main

import (
    "net/http"

    "github.com/Alkemic/go-route/middleware"
)

header := map[string]string{
    "Content-Type": "application/json",
}

middleware.SetHeaders(logger)(
    func(w ResponseWriter, r *Request, p map[string]string) {
        // actual code
    }
)
```

## Examples

### Example function

```go
package main

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
