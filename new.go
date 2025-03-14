package nchi

import (
	"github.com/julienschmidt/httprouter"
	"github.com/muir/nject/v2"
)

type Option func(*rtr)

// rtr is defined the way it is so to prevent users from defining their own Options
type rtr struct {
	*httprouter.Router
}

// The following comment is copied from https://github.com/julienschmidt/httprouter

// WithRedirectTrailingSlash enables/disables automatic redirection if the current route can't be matched but a
// handler for the path with (without) the trailing slash exists.
// For example if /foo/ is requested but a route only exists for /foo, the
// client is redirected to /foo with http status code 301 for GET requests
// and 307 for all other request methods.  The default is: enabled.
func WithRedirectTrailingSlash(b bool) Option {
	return func(r *rtr) {
		r.RedirectTrailingSlash = b
	}
}

// The following comment is copied from https://github.com/julienschmidt/httprouter

// WithRedirectFixedPath enables/disables trying to fix the current request path, if no
// handle is registered for it.
// First superfluous path elements like ../ or // are removed.
// Afterwards the router does a case-insensitive lookup of the cleaned path.
// If a handle can be found for this route, the router makes a redirection
// to the corrected path with status code 301 for GET requests and 307 for
// all other request methods.
// For example /FOO and /..//Foo could be redirected to /foo.
// RedirectTrailingSlash is independent of this option.  The default is: enabled.
func WithRedirectFixedPath(b bool) Option {
	return func(r *rtr) {
		r.RedirectFixedPath = b
	}
}

// The following comment is copied from https://github.com/julienschmidt/httprouter

// WithHandleMethodNotAllowed enables/disables, checking if another method is allowed for the
// current route, if the current request can not be routed.
// If this is the case, the request is answered with 'Method Not Allowed'
// and HTTP status code 405.
// If no other Method is allowed, the request is delegated to the NotFound
// handler.  The default is: enabled..
func WithHandleMethodNotAllowed(b bool) Option {
	return func(r *rtr) {
		r.HandleMethodNotAllowed = b
	}
}

// The following comment is copied from https://github.com/julienschmidt/httprouter

// WIthHandleOPTIONS enables/disables automatic replies to OPTIONS requests.
// Custom OPTIONS handlers take priority over automatic replies.
// The default is: enabled.
func WithHandleOPTIONS(b bool) Option {
	return func(r *rtr) {
		r.HandleOPTIONS = b
	}
}

func NewRouter(options ...Option) *Mux {
	return &Mux{
		providers: nject.Sequence("router"),
		options:   options,
	}
}
