package nchi_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/muir/nchi"

	"github.com/stretchr/testify/assert"
)

func makeDown(s string) func(string) string {
	return func(i string) string {
		return i + s
	}
}

func bottom(s string, w http.ResponseWriter) {
	_, _ = w.Write([]byte(s))
}

type testCase struct {
	path string
	want string
}

func doTest(t *testing.T, mux *nchi.Mux, cases []testCase) {
	doTestMethod(t, mux, "GET", cases)
}

func doTestMethod(t *testing.T, mux *nchi.Mux, method string, cases []testCase) {
	for _, tc := range cases {
		tc := tc
		t.Run(method+" "+tc.path, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(method, tc.path, nil)
			mux.ServeHTTP(w, r)
			body, err := io.ReadAll(w.Result().Body)
			assert.NoError(t, err, tc.path)
			got := string(body)
			t.Log("->", got)
			assert.Equal(t, tc.want, got, tc.path)
		})
	}
}

func TestUse(t *testing.T) {
	mux := nchi.NewRouter()
	mux.Use("")
	mux.Use(makeDown("a"))
	mux.Get("/simplest", bottom)
	mux.Get("/two", makeDown("b"), bottom)

	mux.Use(makeDown("c"), makeDown("d"))
	mux.Get("/afteruse", bottom)

	doTest(t, mux, []testCase{
		{path: "/simplest", want: "a"},
		{path: "/two", want: "ab"},
		{path: "/afteruse", want: "acd"},
	})
}

func TestWith(t *testing.T) {
	mux := nchi.NewRouter()
	mux.Use("")
	mux.Use(makeDown("a"))

	m2 := mux.With(makeDown("b"), makeDown("c"))

	m3 := mux.With(makeDown("d"))
	m3.Use(makeDown("e"))

	mux.Get("/mux", bottom)
	m2.Get("/m2", bottom)
	m3.Get("/m3", bottom)

	doTest(t, mux, []testCase{
		{path: "/mux", want: "a"},
		{path: "/m2", want: "abc"},
		{path: "/m3", want: "ade"},
	})
}

func TestRoute(t *testing.T) {
	mux := nchi.NewRouter()
	mux.Use("")
	mux.Use(makeDown("a"))

	mux.Route("/r2", func(m2 *nchi.Mux) {
		m2.Use(makeDown("b"))
		m2.Get("", makeDown("c"), bottom)
		m2.Get("/a", makeDown("d"), bottom)
		m2.Get("/b", makeDown("e"), bottom)
	})

	mux.Get("/c", makeDown("f"), bottom)

	doTest(t, mux, []testCase{
		{path: "/r2", want: "abc"},
		{path: "/r2/a", want: "abd"},
		{path: "/r2/b", want: "abe"},
		{path: "/c", want: "af"},
	})
}

func TestGroup(t *testing.T) {
	mux := nchi.NewRouter()
	mux.Use("")
	mux.Use(makeDown("a"))

	mux.Group(func(g3 *nchi.Mux) {
		g3.Use("", makeDown("b"))
		g3.Get("/g", bottom)
	})

	mux.Route("/r3", func(m2 *nchi.Mux) {
		m2.Use(makeDown("c"))
		m2.Get("/c", bottom)
		m2.Group(func(g4 *nchi.Mux) {
			g4.Use("", makeDown("d"))
			g4.Get("/d", bottom)
		})
	})

	mux.Get("/e", bottom)

	doTest(t, mux, []testCase{
		{path: "/g", want: "b"},
		{path: "/r3/c", want: "ac"},
		{path: "/r3/d", want: "d"},
		{path: "/e", want: "a"},
	})
}

func TestEndpoint(t *testing.T) {
	mux := nchi.NewRouter()
	mux.Get("/thing/:thingID", func(endpoint nchi.Endpoint, w http.ResponseWriter) {
		_, _ = w.Write([]byte(endpoint))
	})
	w := httptest.NewRecorder()

	r := httptest.NewRequest("GET", "/thing/473", nil)
	mux.ServeHTTP(w, r)
	body, err := io.ReadAll(w.Result().Body)
	assert.NoError(t, err)
	got := string(body)
	t.Log("->", got)
	assert.Equal(t, "/thing/:thingID", got)
}
