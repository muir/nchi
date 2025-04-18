package nchi

import (
	"net/http"

	"github.com/muir/nject/v2"
	"github.com/pkg/errors"

	"github.com/julienschmidt/httprouter"
)

type Params = httprouter.Params

// Endpoint is a type that handlers can accepte as an input.  It will be the
// combined URL path without path variables substituted.  If you have
//
//	mux.Get("/thing/:thingID", handler)
//
// and handler takes an nchi.Endpoint argument, and there is a request for
// http://example.com/thing/3802, then the nchi.Endpoint will be "/thing/:thingID".
type Endpoint string

type Mux struct {
	providers *nject.Collection // partial set
	routes    []*Mux
	path      string             // a fragment
	method    string             // set for endpoints only
	router    *httprouter.Router // only set at top-most
	group     bool
	options   []Option
	special   *special
}

func (mux *Mux) add(n *Mux) *Mux {
	mux.routes = append(mux.routes, n)
	if !n.group {
		n.providers = mux.providers.Append(n.path, n.providers)
	}
	return n
}

// With is just like Use except that it returns a new Mux instead of
// modifying the current one
func (mux *Mux) With(providers ...interface{}) *Mux {
	return mux.add(&Mux{
		providers: nject.Sequence(mux.path, translateMiddleware(providers)...),
	})
}

// Route establishes a new Mux at a new path (combined with the
// current path context).
func (mux *Mux) Route(path string, f func(mux *Mux)) {
	f(mux.add(&Mux{
		path:      path,
		providers: nject.Sequence(mux.path),
	}))
}

// Group establishes a new Mux at the current path but does
// not inherit any middlewhere.
func (mux *Mux) Group(f func(mux *Mux)) {
	f(mux.add(&Mux{
		group:     true,
		providers: nject.Sequence(mux.path),
	}))
}

// Method registers an endpoint handler at the new path (combined with
// the current path) using a combination of inherited middleware and
// the providers here.
func (mux *Mux) Method(method string, path string, providers ...interface{}) {
	mux.add(&Mux{
		providers: nject.Sequence(method+" "+path, translateMiddleware(providers)...),
		method:    method,
		path:      path,
	})
}

func (mux *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if mux.router == nil {
		err := mux.Bind()
		if err != nil {
			panic(err.Error())
		}
	}
	mux.router.ServeHTTP(w, r)
}

// Bind validates that the injection chains for all routes are valid.
// If any are not, an error is returned.  If you do not call bind, and
// there are any invalid injection chains, then routes will panic when
// used.
func (mux *Mux) Bind() error {
	router := httprouter.New()
	for _, opt := range mux.options {
		opt(&rtr{router})
	}
	err := mux.bind(router, "")
	if err != nil {
		return err
	}
	mux.router = router
	return nil
}

func (mux *Mux) bind(router *httprouter.Router, path string) error {
	combinedPath := path + mux.path
	for _, route := range mux.routes {
		err := route.bind(router, combinedPath)
		if err != nil {
			return err
		}
	}
	providers := nject.Sequence(path,
		Endpoint(combinedPath),
		mux.providers,
	)
	if mux.special != nil {
		return mux.bindSpecial(router, combinedPath, providers)
	}
	if mux.method == "" {
		return nil
	}
	var handle httprouter.Handle
	err := providers.Bind(&handle, nil)
	if err != nil {
		return errors.Wrapf(err, "bind router %s %s", mux.method, combinedPath)
	}
	router.Handle(mux.method, combinedPath, handle)
	return nil
}

// Use adds additional http middleware (implementing the http.Handler interface)
// or nject-style providers to the current handler context.  These middleware
// and providers will be injected into the handler chain for any downstream
// endpoints that are defined after the call to Use.
func (mux *Mux) Use(providers ...interface{}) {
	n := "router"
	if mux.path != "" {
		n = mux.path
	}
	mux.providers = mux.providers.Append(n, translateMiddleware(providers)...)
}

// Get establish a route for HTTP GET requests
func (mux *Mux) Get(path string, providers ...interface{}) { mux.Method("GET", path, providers...) }

// Head establish a route for HTTP HEAD requests
func (mux *Mux) Head(path string, providers ...interface{}) { mux.Method("HEAD", path, providers...) }

// Post establish a route for HTTP POST requests
func (mux *Mux) Post(path string, providers ...interface{}) { mux.Method("POST", path, providers...) }

// Put establish a route for HTTP PUT requests
func (mux *Mux) Put(path string, providers ...interface{}) { mux.Method("PUT", path, providers...) }

// Patch establish a route for HTTP PATCH requests
func (mux *Mux) Patch(path string, providers ...interface{}) { mux.Method("PATCH", path, providers...) }

// Options establish a route for HTTP OPTIONS requests.
func (mux *Mux) Options(path string, providers ...interface{}) {
	mux.Method("OPTIONS", path, providers...)
}

// Delete establish a route for HTTP DELETE requests.
func (mux *Mux) Delete(path string, providers ...interface{}) {
	mux.Method("DELETE", path, providers...)
}
