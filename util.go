package lars

import "net/http"

// NativeChainHandler is used in native handler chains
// example using nosurf crsf middleware nosurf.NewPure(lars.NativeChainHandlerFunc)
var NativeChainHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	c := GetContext(w)

	if c.index+1 < len(c.handlers) {
		c.Next()
	}
})

// wrapHandler wraps Handler type
func wrapHandler(h Handler) HandlerFunc {

	switch h := h.(type) {
	case HandlerFunc:
		return h
	case func(*Context):
		return h
	case http.Handler, http.HandlerFunc:
		return func(c *Context) {

			if h.(http.Handler).ServeHTTP(c.Response, c.Request); c.Response.status != http.StatusOK || c.Response.committed {
				return
			}

			if c.index+1 < len(c.handlers) {
				c.Next()
			}
		}
	case func(http.ResponseWriter, *http.Request):
		return func(c *Context) {

			if h(c.Response, c.Request); c.Response.status != http.StatusOK || c.Response.committed {
				return
			}

			if c.index+1 < len(c.handlers) {
				c.Next()
			}
		}
	case func(handler http.Handler) http.Handler:

		hf := h(NativeChainHandler)

		return func(c *Context) {
			hf.ServeHTTP(c.Response, c.Request)
		}
	default:
		panic("unknown handler")
	}
}

// GetContext is a helper method for retrieving the *Context object from
// the ResponseWriter when using native go hanlders.
// NOTE: this will panic if fed an http.ResponseWriter not provided by lars's
// chaining.
func GetContext(w http.ResponseWriter) *Context {
	return w.(*Response).context
}
