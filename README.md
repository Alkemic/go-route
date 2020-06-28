# go-route - a regexp based router

``go-route`` allows you to create multi level regexp-based routing for your application with registering named 
parameters in regexp.

Parameters from regexp routing are passed through request's context. You can fetch parameters using `rotue.GetParams`
function, it accepts only one parameter which is request, and returns `map[string]string`. In returned map keys are named
groups and values are matched values from URL.

## Usage

```go
package main

import (
    ...
    "github.com/Alkemic/go-route"
)

func main() {
    newsRoutes := route.New()
    newsRoutes.Add(`^/$`, view.NewsList, http.MethodGet)
    newsRoutes.Add(`^/(?P<pk>\d+)/$`, view.News, http.MethodGet, http.MethodPost)
    newsRoutes.Add(`^/(?P<pk>\d+),(?P<slug>[a-z\-_]+)\.html$`, view.News, http.MethodGet, http.MethodPost)

    routing := route.New()
    routing.Add(`^/news`, newsRoutes) // register sub routes

    // register custom 404 handler function
    routing.NotFound = helper.Handle404

    log.Fatalln(http.ListenAndServe("0.0.0.0:80", routing))
}
```

To start using `go-route` first instance of `RegexpRouter` must be created using `route.New()` function, then by calling
`Add` method create new route by passing a valid URL regexp, http handler, and optionally methods that can are allowed
for given route. A http handler can be normal http function `func(w http.ResponseWriter, r *http.Request)` or other
`RegexpRouter` instance like in given example.

Keep in mind that when using sub-routing, URL path that will be tested against regexp will be strip from matching
beginning, i.e.: path `/news/123,important-news.html` will be passed to sub router as `/123,important-news.html`. 

Then the base routing should be passed into `http.ListenAndServe`.

### Getting parameters in HTTP handler function

Parameters are pas

```go
package main

import (
    ...
    "github.com/Alkemic/go-route"
)

func News(w http.ResponseWriter, r *http.Request) {
    p := route.GetParams(req)
    pk, _ := strconv.Atoi(p["pk"])
    news, _ := model.GetNews(pk)
    comments, _ := news.GetComments()
    _ = json.NewEncoder(w).Encode(map[string]interface{}{
        "Entry":    news,
        "Comments": comments,
    })
}

```

## Middlewares

### Allowed methods

This middleware checks if request method is allowed.

Use consts from ``net/http`` for methods names, i.e. ``http.MethodPost``.

```go
package main

import (
    "net/http"

    "github.com/Alkemic/go-route/middleware"
)

middleware.AllowedMethods([]string{http.MethodPost, http.MethodPut}])(
    func(w http.ResponseWriter, r *http.Request) {
        p := route.GetParams(req)
        // actual code
    }
)
```

### Panic interceptors

This middleware catches panics and logs them, and then returns 500 to the user. 500 handler can be overridden by

Differences between ``PanicInterceptor`` and ``PanicInterceptorWithLogger`` is that, the second one is parameterized
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
    func(w http.ResponseWriter, r *http.Request) {
        p := route.GetParams(req)
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

    "github.com/Alkemic/go-route"
    "github.com/Alkemic/go-route/middleware"
)

var logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

middleware.PanicInterceptorWithLogger(logger)(
    func(w http.ResponseWriter, r *http.Request) {
        p := route.GetParams(req)
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

    "github.com/Alkemic/go-route"
    "github.com/Alkemic/go-route/middleware"
)

header := map[string]string{
    "Content-Type": "application/json",
}

middleware.SetHeaders(logger)(
    func(w http.ResponseWriter, r *http.Request) {
        p := route.GetParams(req)
        // actual code
    }
)
```

### BasicAuthenticate

Setups HTP basic authenticate, for given view. Accepts functions that will verify credentials. Comes with two default
functions that will do that:
* ``Authenticate`` - accepts user and password
* ``AuthenticateList`` - accept map (user => password)

Custom verification function must have following signature ``func(user string, password string) (string, error)``.

```go
package main

import (
    "net/http"

    "github.com/Alkemic/go-route"
    "github.com/Alkemic/go-route/middleware"
)

authenticateFunc = middleware.Authenticate("username", "password")

middleware.BasicAuthenticate(logger, authenticateFunc, "realm name")(
    func(w http.ResponseWriter, r *http.Request) {
        user, err := middleware.GetUser(r)
    }
)
```

### Noop

Does noting. Simply returns provided functions. Can be useful when used as default option in some cases.
