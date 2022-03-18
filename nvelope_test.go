package nchi_test

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/muir/nchi"
	"github.com/muir/nvelope"

	"github.com/stretchr/testify/assert"
)

type bodyData struct {
	I int `json:"i"`
}

type parameters struct {
	Body bodyData `nvelope:"model"`
	Zoom int      `nvelope:"path,name=zoomie"`
	Who  string   `nvelope:"query,name=who"`
}

func TestDecoders(t *testing.T) {
	mux := nchi.NewRouter(nchi.WithRedirectFixedPath(true))
	mux.Use(nvelope.MinimalErrorHandler, nvelope.ReadBody)
	mux.Post("/td1/:zoomie/foo", nchi.DecodeJSON, func(p parameters, w http.ResponseWriter) {
		_, _ = w.Write([]byte(fmt.Sprintf("bi %d z %d w %s", p.Body.I, p.Zoom, p.Who)))
	})
	mux.Post("/td2/:zoomie/bar", nchi.DecodeXML, func(p parameters, w http.ResponseWriter) {
		_, _ = w.Write([]byte(fmt.Sprintf("bi %d z %d w %s", p.Body.I, p.Zoom, p.Who)))
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/td1/9/foo?who=john", strings.NewReader(`{"I":3}`))
	mux.ServeHTTP(w, r)
	body, err := io.ReadAll(w.Result().Body)
	assert.NoError(t, err, "/td1")
	got := string(body)
	t.Log("/td1/9/foo?who=john ->", got)
	assert.Equal(t, "bi 3 z 9 w john", got, "/td1")

	enc, err := xml.Marshal(bodyData{I: 19})
	assert.NoError(t, err, "marshal xml")
	w = httptest.NewRecorder()
	r = httptest.NewRequest("POST", "/td2/11/bar?who=mary", bytes.NewReader(enc))
	mux.ServeHTTP(w, r)
	body, err = io.ReadAll(w.Result().Body)
	assert.NoError(t, err, "/td2")
	got = string(body)
	t.Log("/td2/11/bar?who=mary ->", got)
	assert.Equal(t, "bi 19 z 11 w mary", got, "/td2")
}
