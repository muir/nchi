package nchi_test

import (
	"testing"

	"github.com/muir/nchi"
)

func TestWithRedirectTrailingSlash(t *testing.T) {
	mux := nchi.NewRouter(nchi.WithRedirectTrailingSlash(true))
	mux.Use("")
	mux.Use(makeDown("a"))
	mux.Get("/ts1", makeDown("b"), bottom)

	doTest(t, mux, []testCase{
		{path: "/ts1", want: "ab"},
		{path: "/ts1/", want: "<a href=\"/ts1\">Moved Permanently</a>.\n\n"},
	})

	mux = nchi.NewRouter(nchi.WithRedirectTrailingSlash(false))
	mux.Use("")
	mux.Use(makeDown("a"))
	mux.Get("/ts2", makeDown("b"), bottom)

	doTest(t, mux, []testCase{
		{path: "/ts2", want: "ab"},
		{path: "/ts2/", want: "404 page not found\n"},
	})

	mux = nchi.NewRouter()
	mux.Use("")
	mux.Use(makeDown("a"))
	mux.Get("/ts3", makeDown("b"), bottom)

	doTest(t, mux, []testCase{
		{path: "/ts3", want: "ab"},
		{path: "/ts3/", want: "<a href=\"/ts3\">Moved Permanently</a>.\n\n"},
	})
}

func TestWithRedirectFixedPath(t *testing.T) {
	mux := nchi.NewRouter(nchi.WithRedirectFixedPath(true))
	mux.Use("")
	mux.Use(makeDown("a"))
	mux.Get("/ts1", makeDown("b"), bottom)

	doTest(t, mux, []testCase{
		{path: "/ts1", want: "ab"},
		{path: "/xx/..//ts1/", want: "<a href=\"/ts1\">Moved Permanently</a>.\n\n"},
	})

	mux = nchi.NewRouter(nchi.WithRedirectFixedPath(false))
	mux.Use("")
	mux.Use(makeDown("a"))
	mux.Get("/ts2", makeDown("b"), bottom)

	doTest(t, mux, []testCase{
		{path: "/ts2", want: "ab"},
		{path: "/xx/..//ts2/", want: "404 page not found\n"},
	})

	mux = nchi.NewRouter()
	mux.Use("")
	mux.Use(makeDown("a"))
	mux.Get("/ts3", makeDown("b"), bottom)

	doTest(t, mux, []testCase{
		{path: "/ts3", want: "ab"},
		{path: "/xx/..//ts3", want: "<a href=\"/ts3\">Moved Permanently</a>.\n\n"},
	})
}

func TestWithHandleMethodNotAllowed(t *testing.T) {
	mux := nchi.NewRouter(nchi.WithHandleMethodNotAllowed(true))
	mux.Use("")
	mux.Use(makeDown("a"))
	mux.Post("/mna1", makeDown("b"), bottom)

	doTest(t, mux, []testCase{
		{path: "/mna1", want: "Method Not Allowed\n"},
	})

	mux = nchi.NewRouter(nchi.WithHandleMethodNotAllowed(false))
	mux.Use("")
	mux.Use(makeDown("a"))
	mux.Post("/mna2", makeDown("b"), bottom)

	doTest(t, mux, []testCase{
		{path: "/mna2", want: "404 page not found\n"},
	})

	mux = nchi.NewRouter()
	mux.Use("")
	mux.Use(makeDown("a"))
	mux.Post("/mna3", makeDown("b"), bottom)

	doTest(t, mux, []testCase{
		{path: "/mna3", want: "Method Not Allowed\n"},
	})
}

func TestWithHandleOPTIONS(t *testing.T) {
	mux := nchi.NewRouter(nchi.WithHandleOPTIONS(true))
	mux.Use("")
	mux.Use(makeDown("a"))
	mux.Get("/ho1", makeDown("b"), bottom)
	mux.Options("/ho1", makeDown("c"), bottom)
	mux.Get("/ho2", makeDown("c"), bottom)

	doTestMethod(t, mux, "OPTIONS", []testCase{
		{path: "/ho1", want: "ac"},
		{path: "/ho2", want: ""},
		{path: "/ho3", want: "404 page not found\n"},
	})

	mux = nchi.NewRouter(nchi.WithHandleOPTIONS(false))
	mux.Use("")
	mux.Use(makeDown("a"))
	mux.Get("/ho3", makeDown("b"), bottom)
	mux.Options("/ho3", makeDown("c"), bottom)
	mux.Get("/ho4", makeDown("c"), bottom)

	doTestMethod(t, mux, "OPTIONS", []testCase{
		{path: "/ho3", want: "ac"},
		{path: "/ho4", want: "Method Not Allowed\n"},
		{path: "/ho5", want: "404 page not found\n"},
	})

	mux = nchi.NewRouter()
	mux.Use("")
	mux.Use(makeDown("a"))
	mux.Get("/ho6", makeDown("b"), bottom)
	mux.Options("/ho6", makeDown("c"), bottom)
	mux.Get("/ho7", makeDown("c"), bottom)

	doTestMethod(t, mux, "OPTIONS", []testCase{
		{path: "/ho6", want: "ac"},
		{path: "/ho7", want: ""},
		{path: "/ho8", want: "404 page not found\n"},
	})
}
