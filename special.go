package nchi

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/muir/nject"

	"github.com/pkg/errors"
)

// only one field of special will be set at a time
type special struct {
	serveFiles       http.FileSystem
	globalOPTIONS    bool
	methodNotAllowed bool
	notFound         bool
	panicHandler     bool
}

func (mux *Mux) bindSpecial(router *httprouter.Router, combinedPath string, combinedProviders *nject.Collection) error {
	switch {
	case mux.special.serveFiles != nil:
		router.ServeFiles(combinedPath, mux.special.serveFiles)
	case mux.special.globalOPTIONS:
		if router.GlobalOPTIONS == nil {
			return errors.Wrap(combinedProviders.Bind(&router.GlobalOPTIONS, nil), "bind GlobalOPTIONS")
		}
	case mux.special.methodNotAllowed:
		if router.MethodNotAllowed == nil {
			return errors.Wrap(combinedProviders.Bind(&router.MethodNotAllowed, nil), "bind MethodNotAllowed")
		}
	case mux.special.notFound:
		if router.NotFound == nil {
			return errors.Wrap(combinedProviders.Bind(&router.NotFound, nil), "bind NotFound")
		}
	case mux.special.panicHandler:
		if router.PanicHandler == nil {
			var ph func(w http.ResponseWriter, r *http.Request, rec RecoverInterface)
			err := combinedProviders.Bind(&ph, nil)
			if err != nil {
				return errors.Wrap(err, "bind PanicHandler")
			}
			router.PanicHandler = func(w http.ResponseWriter, r *http.Request, rec interface{}) {
				ph(w, r, rec)
			}
		}
	default:
		return errors.New("should not get here")
	}
	return nil
}

func (mux *Mux) addSpecial(name string, providers []interface{}) *Mux {
	return mux.add(&Mux{
		providers: nject.Sequence(name, translateMiddleware(providers)...),
		special:   &special{},
	})
}

// The following comment is derrived from https://github.com/julienschmidt/httprouter

// GlobalOPTIONS sets a handler that is called on automatically on OPTIONS requests.
// The handler is only called if HandleOPTIONS is true and no OPTIONS
// handler for the specific path was set.
// The "Allowed" header is set before calling the handler.
//
// Only the first GobalOPTIONS call counts.
func (mux *Mux) GlobalOPTIONS(providers ...interface{}) {
	mux.addSpecial("globalOPTIONS", providers).special.globalOPTIONS = true
}

// The following comment is derrived from https://github.com/julienschmidt/httprouter

// NotFound sets a handler
// which is called when no matching route is
// found. If it is not set, http.NotFound is used.
//
// Only the first NotFound call counts.
func (mux *Mux) NotFound(providers ...interface{}) {
	mux.addSpecial("notFound", providers).special.notFound = true
}

// The following comment is derrived from https://github.com/julienschmidt/httprouter

// MethodNotAllowed sets a handler
// which is called when a request
// cannot be routed and HandleMethodNotAllowed is true.
// If it is not set, http.Error with http.StatusMethodNotAllowed is used.
// The "Allow" header with allowed request methods is set before the handler
// is called.
//
// Only the first MethodNotFound call counts.
func (mux *Mux) MethodNotAllowed(providers ...interface{}) {
	mux.addSpecial("methodNotAllowed", providers).special.methodNotAllowed = true
}

type RecoverInterface interface{}

// The following comment is derrived from https://github.com/julienschmidt/httprouter

// PanicHandler sets a handler
// to handle panics recovered from http handlers.
// It should be used to generate a error page and return the http error code
// 500 (Internal Server Error).
// The handler can be used to keep your server from crashing because of
// unrecovered panics.
//
// The type RecoverInterface can be used to receive the interface{} that
// is returned from recover().
//
// Alternatively, use the nvelope.CatchPanic middleware to catch panics.
//
// Only the first PanicHandler call counts.
func (mux *Mux) PanicHandler(providers ...interface{}) {
	mux.addSpecial("panicHandler", providers).special.panicHandler = true
}

// The following comment is copied from https://github.com/julienschmidt/httprouter

// ServeFiles serves files from the given file system root. The path
// must end with "/*filepath", files are then served from the local path
// /defined/root/dir/*filepath. For example if root is "/etc" and *filepath
// is "passwd", the local file "/etc/passwd" would be served. Internally
// a http.FileServer is used, therefore http.NotFound is used instead of
// the Router's NotFound handler. To use the operating system's file system
// implementation, use http.Dir:
//
// Currently, ServeFiles does not use any middleware.  That may change in
// a future release.
func (mux *Mux) ServeFiles(path string, fs http.FileSystem) {
	mux.add(&Mux{
		path: path,
		special: &special{
			serveFiles: fs,
		},
	})
}
