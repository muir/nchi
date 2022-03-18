package nchi_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/muir/nchi"
)

func TestGlobalOPTIONS(t *testing.T) {
	mux := nchi.NewRouter()
	mux.Use("")
	mux.Use(makeDown("a"))
	mux.Options("/go1", makeDown("b"), bottom)
	mux.GlobalOPTIONS(makeDown("c"), bottom)
	mux.Get("/go2", makeDown("d"), bottom)

	doTestMethod(t, mux, "OPTIONS", []testCase{
		{path: "/go1", want: "ab"},
		{path: "/go2", want: "ac"},
		{path: "/go3", want: "404 page not found\n"},
	})
}

func TestNotFound(t *testing.T) {
	mux := nchi.NewRouter()
	mux.Use("")
	mux.Use(makeDown("a"))
	mux.Get("/nf1", makeDown("b"), bottom)
	mux.NotFound(makeDown("c"), bottom)

	doTest(t, mux, []testCase{
		{path: "/nf1", want: "ab"},
		{path: "/nf2", want: "ac"},
	})
}

func TestMethodNotAllowed(t *testing.T) {
	mux := nchi.NewRouter()
	mux.Use("")
	mux.Use(makeDown("a"))
	mux.Post("/mna1", makeDown("b"), bottom)
	mux.MethodNotAllowed(makeDown("c"), bottom)

	doTest(t, mux, []testCase{
		{path: "/mna1", want: "ac"},
		{path: "/mna2", want: "404 page not found\n"},
	})
}

func TestPanicHandler(t *testing.T) {
	mux := nchi.NewRouter()
	mux.Use("")
	mux.Use(makeUp1("a"))
	mux.Get("/ph2", func(w http.ResponseWriter) {
		x := make([]int, 3)
		_, _ = w.Write([]byte(fmt.Sprint(x[5])))
	})
	mux.PanicHandler(func(w http.ResponseWriter, r nchi.RecoverInterface) {
		_, _ = w.Write([]byte(fmt.Sprintf("recover %T-", r)))
	})
	mux.Get("/ph1", bottom)

	doTest(t, mux, []testCase{
		{path: "/ph1", want: "a"},
		{path: "/ph2", want: "recover runtime.boundsError-a"},
	})
}
