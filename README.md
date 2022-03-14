# nchi - http router

[![GoDoc](https://godoc.org/github.com/muir/nchi?status.png)](https://pkg.go.dev/github.com/muir/nchi)
![unit tests](https://github.com/muir/nchi/actions/workflows/go.yml/badge.svg)
[![report card](https://goreportcard.com/badge/github.com/muir/nchi)](https://goreportcard.com/report/github.com/muir/nchi)
[![codecov](https://codecov.io/gh/muir/nchi/branch/main/graph/badge.svg)](https://codecov.io/gh/muir/nchi)

<!---
This readme is heavily based upon the README.md in https://github.com/go-chi/chi.
--->

`nchi` is a lightweight, idiomatic and composable router for building Go HTTP services. It's
especially good at helping you write large REST API services that are kept maintainable as your
project grows and changes. `nchi` is built on top of the 
[nject](https://github.com/muir/nject) dependency injection framework and 
the fastest Go http router, [httprouter](https://github.com/julienschmidt/httprouter).
`nchi` is a straight-up rip-off of [chi](https://github.com/go-chi/chi)
substituting nject for context and in the process making it easier to write middleware and
and endpoints.

`nchi` can use standard middleware and it can use dependency-injection middleware.  See:

- [nvelope](https://github.com/muir/nvelope) for nject-based middleware
- [chi](https://pkg.go.dev/github.com/go-chi/chi/middleware) for chi's middleware collection

Note: if you're using `nvelope.DeferredWriter`, avoid other middleware that replaces
the `http.ResponseWriter`.

"Standard" middlewhare has one of the following shapes:

- `func(http.HandlerFunc) http.HandlerFunc`
- `func(http.Handler) http.Handler`

`nchi` automatically detects standard middleware and translates it for use in
an nject-based framework.

## Install

	go get github.com/muir/nchi

## Examples

**As easy as:**

```go
package main

import (
	"net/http"

	"github.com/muir/nchi"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := nchi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	http.ListenAndServe(":3000", r)
}
```

**REST Preview:**

Here is a little preview of how routing looks like with `nchi`. 

```go
import (
  //...
  "github.com/muir/nchi"
  "github.com/muir/nvelope"
  "github.com/muir/nject"
  "github.com/go-chi/chi/v5/middleware"
)

func main() {
  r := nchi.NewRouter()

  // A good base middleware stack
  r.Use(
     middleware.RequestID, 
     middleware.RealIP, 
     middleware.Logger,
     nchi.DecodeJSON,
  )

  r.Use(func(inner func() error, w http.ResponseWriter) {
     err := inner()
     if err == nil { 
       return
     }
     code := nvelope.GetReturnCode(err)
     w.WriteHeader(code)
     w.Write([]byte(err.Error()))
  })

  // Set a timeout value on the request context (ctx), that will signal
  // through ctx.Done() that the request has timed out and further
  // processing should be stopped.
  r.Use(middleware.Timeout(60 * time.Second))

  r.Get("/", func(w http.ResponseWriter) {
    w.Write([]byte("hi"))
  })

  // RESTy routes for "articles" resource
  r.Route("/articles", func(r nchi.Router) {
    r.With(paginate).Get("/", listArticles)                           // GET /articles
    r.With(paginate).Get("/:month/:day/:year", listArticlesByDate)    // GET /articles/01-16-2017

    r.Post("/", createArticle)                                        // POST /articles
    r.Get("/search", searchArticles)                                  // GET /articles/search

    // Regexp url parameters:
    r.Get("/:articleSlug", getArticleBySlug)                          // GET /articles/home-is-toronto

    // Subrouters:
    r.Route("/:articleID", func(r nchi.Router) {
      r.Use(LoadArticle)
      r.Get("/", getArticle)                                          // GET /articles/123
      r.Put("/", updateArticle)                                       // PUT /articles/123
      r.Delete("/", deleteArticle)                                    // DELETE /articles/123
    })
  })

  // Mount the admin sub-router
  r.Mount("/admin", adminRouter())

  http.ListenAndServe(":3333", r)
}

func LoadArticle(params nchi.Params) (*Article, nject.TerminalError) {
  articleID := params.ByName("articleID")
  article, err := dbGetArticle(articleID)
  if errors.Is(err, sql.NotFound) {
    return nil, nvelope.NotFound(err)
  }
  return article, err
}

func getArticle(article *Article, w http.ResponseWriter) {
  w.Write([]byte(fmt.Sprintf("title:%s", article.Title)))
}
```

