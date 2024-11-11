package nchi_test

import (
	"errors"
	"fmt"
	"io"
	"net/http/httptest"
	"strings"

	"github.com/muir/nchi"
	"github.com/muir/nvelope"
)

type PostBodyModel struct {
	Use      string `json:"use"`
	Exported string `json:"exported"`
	Names    string `json:"names"`
}

type ExampleRequestBundle struct {
	Request     PostBodyModel `nvelope:"model"`
	With        *string       `nvelope:"path,name=with"`
	Parameters  int64         `nvelope:"path,name=parameters"`
	Friends     []int         `nvelope:"query,name=friends"`
	ContentType string        `nvelope:"header,name=Content-Type"`
}

type ExampleResponse struct {
	Stuff string `json:"stuff,omitempty"`
	Here  string `json:"here,omitempty"`
}

func HandleExampleEndpoint(req ExampleRequestBundle) (nvelope.Response, error) {
	if req.ContentType != "application/json" {
		return nil, errors.New("content type must be application/json")
	}
	switch req.Parameters {
	case 666:
		panic("something is not right")
	case 100:
		return nil, nil
	default:
		return ExampleResponse{
			Stuff: *req.With,
		}, nil
	}
}

func Service(r *nchi.Mux) {
	r.Use(
		// order matters and this is a correct order
		nvelope.NoLogger,
		nvelope.InjectWriter,
		nvelope.EncodeJSON,
		nvelope.CatchPanic,
		nvelope.Nil204,
		nvelope.ReadBody,
		nchi.DecodeJSON,
	)
	r.Post("/a/path/:with/:parameters",
		HandleExampleEndpoint,
	)
}

// Example shows an injection chain handling a single endpoint using nject,
// nape, and nvelope.
func Example() {
	r := nchi.NewRouter()
	Service(r)
	ts := httptest.NewServer(r)
	client := ts.Client()
	doPost := func(url string, body string) {
		// nolint:noctx
		res, err := client.Post(ts.URL+url, "application/json",
			strings.NewReader(body))
		if err != nil {
			fmt.Println("response error:", err)
			return
		}
		b, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("read error:", err)
			return
		}
		res.Body.Close()
		fmt.Println(res.StatusCode, "->"+string(b))
	}
	doPost("/a/path/joe/37", `{"Use":"yeah","Exported":"uh hu"}`)
	doPost("/a/path/joe/100", `{"Use":"yeah","Exported":"uh hu"}`)
	doPost("/a/path/joe/38", `invalid json`)
	doPost("/a/path/joe/666", `{"Use":"yeah","Exported":"uh hu"}`)

	// Output: 200 ->{"stuff":"joe"}
	// 204 ->
	// 400 ->nchi_test.ExampleRequestBundle model: Could not decode application/json into nchi_test.PostBodyModel: invalid character 'i' looking for beginning of value
	// 500 ->panic: something is not right
}
