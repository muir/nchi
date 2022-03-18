package nchi

import (
	"net/http"

	"github.com/muir/nvelope"
)

func translateMiddleware(raw []interface{}) []interface{} {
	n := make([]interface{}, 0, len(raw))
	hms := make([]func(http.Handler) http.Handler, 0, len(raw))
	hfs := make([]func(http.HandlerFunc) http.HandlerFunc, 0, len(raw))
	for i := 0; i < len(raw); i++ {
		if h, ok := raw[i].(func(http.Handler) http.Handler); ok {
			hms = append(hms, h)
			var j int
			for j = i + 1; j < len(raw); j++ {
				h, ok := raw[j].(func(http.Handler) http.Handler)
				if !ok {
					break
				}
				hms = append(hms, h)
			}
			n = append(n, nvelope.MiddlewareHandlerBaseWriter(hms...))
			i = j - 1
			hms = hms[:0]
		} else if h, ok := raw[i].(func(http.HandlerFunc) http.HandlerFunc); ok {
			hfs = append(hfs, h)
			var j int
			for j = i + 1; j < len(raw); j++ {
				h, ok := raw[j].(func(http.HandlerFunc) http.HandlerFunc)
				if !ok {
					break
				}
				hfs = append(hfs, h)
			}
			n = append(n, nvelope.MiddlewareBaseWriter(hfs...))
			i = j - 1
			hfs = hfs[:0]
		} else {
			n = append(n, raw[i])
		}
	}
	return n
}
