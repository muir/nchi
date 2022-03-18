package nchi_test

import (
	"net/http"
	"testing"

	"github.com/muir/nchi"
)

func makeUp1(s string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			_, _ = w.Write([]byte(s))
		})
	}
}

func makeUp2(s string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			_, _ = w.Write([]byte(s))
		})
	}
}

func TestMiddlewareHandler(t *testing.T) {
	mux := nchi.NewRouter()
	mux.Use("")
	mux.Use(makeDown("a"))
	mux.Get("/up1", makeUp1("b"), bottom)

	mux.Use(makeUp2("c"), makeUp2("d"))
	mux.Get("/up2", bottom)

	mux.Use(makeUp1("e"), makeUp1("f"))
	mux.Get("/up3", bottom)

	doTest(t, mux, []testCase{
		{path: "/up1", want: "ab"},
		{path: "/up2", want: "adc"},
		{path: "/up3", want: "afedc"},
	})
}
