package web

import (
	"fmt"
	"net/http"
)

const unknownHandler = `Unknown handler type %T. See http://godoc.org/github.com/zenazn/goji/web#HandlerType for a list of acceptable types.`

type netHTTPHandlerWrap struct{ http.Handler }
type netHTTPHandlerFuncWrap struct {
	fn func(http.ResponseWriter, *http.Request)
}
type handlerFuncWrap struct {
	fn func(C, http.ResponseWriter, *http.Request)
}

func (h netHTTPHandlerWrap) ServeHTTPC(c C, w http.ResponseWriter, r *http.Request) {
	h.Handler.ServeHTTP(w, r)
}
func (h netHTTPHandlerFuncWrap) ServeHTTPC(c C, w http.ResponseWriter, r *http.Request) {
	h.fn(w, r)
}
func (h handlerFuncWrap) ServeHTTPC(c C, w http.ResponseWriter, r *http.Request) {
	h.fn(c, w, r)
}

func parseHandler(h HandlerType, di Injector) Handler {
	fn := h
check:
	switch f := fn.(type) {
	case func(c C, w http.ResponseWriter, r *http.Request):
		return handlerFuncWrap{f}
	case func(w http.ResponseWriter, r *http.Request):
		return netHTTPHandlerFuncWrap{f}
	case Handler:
		return f
	case http.Handler:
		return netHTTPHandlerWrap{f}
	default:
		if di != nil {
			var err error
			if fn, err = di.Inject(h); err == nil {
				di = nil
				goto check
			}
			panic(fmt.Errorf("error injecting handler (%v): %v", h, err))
		}
		panic(fmt.Sprintf(unknownHandler, h))
	}
}
