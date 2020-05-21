package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var testHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("bar"))
})

func TestCORSHeader(t *testing.T) {
	cases := []struct {
		name       string
		method     string
		reqHeaders map[string]string
		statusCode int
	}{
		{
			"No header",
			"GET",
			map[string]string{},
			http.StatusBadRequest,
		},
		{
			"GET with accepted headers",
			"GET",
			map[string]string{
				"Origin":                       "http://test.com",
				"Access-Control-Allow-Method":  "GET",
				"Access-Control-Allow-Headers": "Content-Type, Authorization",
			},
			http.StatusOK,
		},
		{
			"GET with Content-Type",
			"GET",
			map[string]string{
				"Origin":                       "http://test.com",
				"Access-Control-Allow-Method":  "GET",
				"Access-Control-Allow-Headers": "Content-Type",
			},
			http.StatusOK,
		},
		{
			"GET with Content-Type",
			"GET",
			map[string]string{
				"Origin":                       "http://test.com",
				"Access-Control-Allow-Method":  "GET",
				"Access-Control-Allow-Headers": "Authorization",
			},
			http.StatusOK,
		},
		{
			"PUT with accepted headers",
			"PUT",
			map[string]string{
				"Origin":                       "http://test.com",
				"Access-Control-Allow-Method":  "PUT",
				"Access-Control-Allow-Headers": "Content-Type, Authorization",
			},
			http.StatusOK,
		},
		{
			"POST with accepted headers",
			"POST",
			map[string]string{
				"Origin":                       "http://test.com",
				"Access-Control-Allow-Method":  "POST",
				"Access-Control-Allow-Headers": "Content-Type, Authorization",
			},
			http.StatusOK,
		},
		{
			"PATCH with accepted headers",
			"PATCH",
			map[string]string{
				"Origin":                       "http://test.com",
				"Access-Control-Allow-Method":  "PATCH",
				"Access-Control-Allow-Headers": "Content-Type, Authorization",
			},
			http.StatusOK,
		},
		{
			"DELETE with accepted headers",
			"DELETE",
			map[string]string{
				"Origin":                       "http://test.com",
				"Access-Control-Allow-Method":  "DELETE",
				"Access-Control-Allow-Headers": "Content-Type, Authorization",
			},
			http.StatusOK,
		},
		{
			"Accepted headers, different order",
			"GET",
			map[string]string{
				"Origin":                       "http://test.com",
				"Access-Control-Allow-Method":  "GET",
				"Access-Control-Allow-Headers": "Authorization, Content-Type",
			},
			http.StatusOK,
		},
		{
			"Unaccepted headers",
			"GET",
			map[string]string{
				"Origin":                       "http://test.com",
				"Access-Control-Allow-Method":  "GET",
				"Access-Control-Allow-Headers": "X-Header-2, X-HEADER-1",
			},
			http.StatusBadRequest,
		},
		{
			"Different Website",
			"GET",
			map[string]string{
				"Origin":                        "http://differenttest.com/",
				"Access-Control-Request-Method": "GET",
				"Access-Control-Allow-Headers":  "Content-Type, Authorization",
			},
			http.StatusOK,
		},
		{
			"Unaccepted Request Method",
			"CONNECT",
			map[string]string{
				"Origin":                        "http://test.com/",
				"Access-Control-Request-Method": "CONNECT",
				"Access-Control-Allow-Headers":  "Content-Type, Authorization",
			},
			http.StatusBadRequest,
		},
		{
			"Accepted and Unaccepted Request Headers",
			"CONNECT",
			map[string]string{
				"Origin":                        "http://test.com/",
				"Access-Control-Request-Method": "CONNECT",
				"Access-Control-Allow-Headers":  "Content-Type, X-Header-2",
			},
			http.StatusBadRequest,
		},
	}

	for _, c := range cases {
		resp := httptest.NewRecorder()
		req, _ := http.NewRequest("OPTIONS", "http://example.com", nil)
		for name, value := range c.reqHeaders {
			req.Header.Add(name, value)
		}
		NewCORSHeader(testHandler).ServeHTTP(resp, req)
		allowedMethods := resp.Header().Get("Access-Control-Allow-Methods")
		allowedOrigin := resp.Header().Get("Access-Control-Allow-Origin")
		allowedHeaders := resp.Header().Get("Access-Control-Allow-Headers")
		exposeHeaders := resp.Header().Get("Access-Control-Expose-Headers")
		maxAge := resp.Header().Get("Access-Control-Max-Age")
		if allowedMethods != "GET, PUT, POST, PATCH, DELETE" {
			t.Errorf("case: %s\nexpected value %s but got %s", c.name, "GET, PUT, POST, PATCH, DELETE", allowedMethods)
		}
		if allowedOrigin != "*" {
			t.Errorf("case: %s\nexpected value %s but got %s", c.name, "*", allowedOrigin)
		}
		if allowedHeaders != "Content-Type, Authorization" {
			t.Errorf("case: %s\nexpected value %s but got %s", c.name, "Content-Type, Authorization", allowedHeaders)
		}
		if exposeHeaders != "Authorization" {
			t.Errorf("case: %s\nexpected value %s but got %s", c.name, "Authorization", exposeHeaders)
		}
		if maxAge != "600" {
			t.Errorf("case: %s\nexpected value %s but got %s", c.name, "600", maxAge)
		}
	}
}
