package handlers

import (
	"net/http"
)

//CORSHeader is a middleware handler that adds a CORS Header
type CORSHeader struct {
	handler http.Handler
}

//ServeHTTP adds the CORS header before passing the repsonse and request
//to the real handler
func (c *CORSHeader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, PATCH, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Expose-Headers", "Authorization")
	w.Header().Set("Access-Control-Max-Age", "600")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	c.handler.ServeHTTP(w, r)
}

//NewCORSHeader constructs a new CORSHeader middlware handler
func NewCORSHeader(handlerToWrap http.Handler) *CORSHeader {
	return &CORSHeader{handlerToWrap}
}
