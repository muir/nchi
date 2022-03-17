package nchi

import (
	"net/http"

	"github.com/muir/nject"
	"github.com/pkg/errors"

	"github.com/julienschmidt/httprouter"
)

type Params = httprouter.Params

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

// With is just like Use except that it returns a new Mux instead of
// modifying the current one
func (mux *Mux) With(providers ...interface{}) *Mux {
	n := &Mux{
		providers: nject.Sequence(mux.path, translateMiddleware(providers)...),
	}
	return n
}

/*
func (mux *Mux) add(n *Mux) *Mux{
	mux.routes = append(mux.routes, n)
	if !n.group {

*/

// Route establishes a new Mux at a new path (combined with the
// current path context).
func (mux *Mux) Route(path string, f func(mux *Mux)) {
	n := &Mux{
		providers: nject.Sequence(mux.path),
	}
	mux.routes = append(mux.routes, n)
	f(n)
}

// Group establishes a new Mux at the current path but does
// not inherit any middlewhere.
func (mux *Mux) Group(f func(mux *Mux)) {
	n := &Mux{
		group:     true,
		providers: nject.Sequence(mux.path),
	}
	mux.routes = append(mux.routes, n)
	f(n)
}

// Method registers an endpoint handler at the new path (combined with
// the current path) using a combination of inherited middleware and
// the providers here.
func (mux *Mux) Method(method string, path string, providers ...interface{}) {
	n := &Mux{
		providers: nject.Sequence(method+" "+path, translateMiddleware(providers)...),
		method:    method,
		path:      path,
	}
	mux.routes = append(mux.routes, n)
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

func (mux *Mux) Bind() error {
	router := httprouter.New()
	for _, opt := range mux.options {
		opt(&rtr{router})
	}
	err := mux.bind(router, "", nject.Sequence("empty"))
	if err != nil {
		return err
	}
	mux.router = router
	return nil
}

func (mux *Mux) bind(router *httprouter.Router, path string, providers *nject.Collection) error {
	combinedPath := path + mux.path
	var combinedProviders *nject.Collection
	if mux.group {
		combinedProviders = mux.providers
	} else {
		combinedProviders = providers.Append(combinedPath, mux.providers)
	}
	for _, route := range mux.routes {
		err := route.bind(router, combinedPath, combinedProviders)
		if err != nil {
			return err
		}
	}
	if mux.special != nil {
		return mux.bindSpecial(router, combinedPath, combinedProviders)
	}
	if mux.method == "" {
		return nil
	}
	var handle httprouter.Handle
	err := combinedProviders.Bind(&handle, nil)
	if err != nil {
		return errors.Wrapf(err, "bind router %s %s", mux.method, combinedPath)
	}
	router.Handle(mux.method, combinedPath, handle)
	return nil
}

// Use adds additional http middleware (implementing the http.Handler interface)
// or nject-style providers to the current handler context.  These middleware
// and providers will be injected into the handler chain for any downstream
// endpoints.
func (mux *Mux) Use(providers ...interface{}) {
	n := "router"
	if mux.path != "" {
		n = mux.path
	}
	mux.providers = mux.providers.Append(n, translateMiddleware(providers)...)
}

func (mux *Mux) Get(path string, providers ...interface{})   { mux.Method("GET", path, providers...) }
func (mux *Mux) Head(path string, providers ...interface{})  { mux.Method("HEAD", path, providers...) }
func (mux *Mux) Post(path string, providers ...interface{})  { mux.Method("POST", path, providers...) }
func (mux *Mux) Put(path string, providers ...interface{})   { mux.Method("PUT", path, providers...) }
func (mux *Mux) Patch(path string, providers ...interface{}) { mux.Method("PATCH", path, providers...) }
func (mux *Mux) Options(path string, providers ...interface{}) {
	mux.Method("OPTIONS", path, providers...)
}

func (mux *Mux) Delete(path string, providers ...interface{}) {
	mux.Method("DELETE", path, providers...)
}
