package mux

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func writeSuccess(w http.ResponseWriter, response string) {
	w.Write([]byte(response))
}

func newRequest(method string, path string) *http.Request {
	req, _ := http.NewRequest(method, path, nil)
	return req
}

func checkSuccess(t *testing.T, record *httptest.ResponseRecorder, response string) {
	if p, err := ioutil.ReadAll(record.Body); err != nil {
		t.Fail()
	} else {
		if !strings.Contains(string(p), response) {
			t.Errorf("should contain '%s'", response)
		}
	}
}

func requestAndCheck(t *testing.T, router *Router, method string, path string, response string) (*http.Request, *httptest.ResponseRecorder) {
	recorder := httptest.NewRecorder()
	req := newRequest(method, path)
	router.ServeHTTP(recorder, req)
	checkSuccess(t, recorder, response)
	return req, recorder
}

/* test basic routes */
func TestRoutes(t *testing.T) {
	router := NewRouter()
	defer router.Free()

	/* gets */
	router.Get("/", func(w http.ResponseWriter, req *http.Request) {
		writeSuccess(w, "get test1")
	})
	router.Get("/foo", func(w http.ResponseWriter, req *http.Request) {
		writeSuccess(w, "get test2")
	})
	router.Get("/foo/{id}", func(w http.ResponseWriter, req *http.Request) {
		writeSuccess(w, "get test3")
	})

	/* posts */
	router.Post("/", func(w http.ResponseWriter, req *http.Request) {
		writeSuccess(w, "post test1")
	})
	router.Post("/foo", func(w http.ResponseWriter, req *http.Request) {
		writeSuccess(w, "post test2")
	})

	/* puts */
	router.Put("/", func(w http.ResponseWriter, req *http.Request) {
		writeSuccess(w, "put test1")
	})
	router.Put("/foo", func(w http.ResponseWriter, req *http.Request) {
		writeSuccess(w, "put test2")
	})

	/* deletes */
	router.Delete("/", func(w http.ResponseWriter, req *http.Request) {
		writeSuccess(w, "delete test1")
	})
	router.Delete("/foo", func(w http.ResponseWriter, req *http.Request) {
		writeSuccess(w, "delete test2")
	})

	/* patches */
	router.Patch("/", func(w http.ResponseWriter, req *http.Request) {
		writeSuccess(w, "patch test1")
	})
	router.Patch("/foo", func(w http.ResponseWriter, req *http.Request) {
		writeSuccess(w, "patch test2")
	})

	/* options */
	router.Options("/", func(w http.ResponseWriter, req *http.Request) {
		writeSuccess(w, "options test1")
	})
	router.Options("/foo", func(w http.ResponseWriter, req *http.Request) {
		writeSuccess(w, "options test2")
	})

	router.Compile()

	requestAndCheck(t, router, "GET", "/", "get test1")
	requestAndCheck(t, router, "GET", "/foo", "get test2")
	requestAndCheck(t, router, "GET", "/foo/teststr", "get test3")
	requestAndCheck(t, router, "POST", "/", "post test1")
	requestAndCheck(t, router, "POST", "/foo", "post test2")
	requestAndCheck(t, router, "PUT", "/", "put test1")
	requestAndCheck(t, router, "PUT", "/foo", "put test2")
	requestAndCheck(t, router, "DELETE", "/", "delete test1")
	requestAndCheck(t, router, "DELETE", "/foo", "delete test2")
	requestAndCheck(t, router, "PATCH", "/", "patch test1")
	requestAndCheck(t, router, "PATCH", "/foo", "patch test2")
	requestAndCheck(t, router, "OPTIONS", "/", "options test1")
	requestAndCheck(t, router, "OPTIONS", "/foo", "options test2")
}

/* testing path variables */
func TestPathVars(t *testing.T) {
	var req *http.Request
	router := NewRouter()
	defer router.Free()

	router.Get("/foo/{id}", func(w http.ResponseWriter, req *http.Request) {
		writeSuccess(w, "get pathtest1")
	})

	router.Get("/foo/{id}/{var2}", func(w http.ResponseWriter, req *http.Request) {
		writeSuccess(w, "get pathtest2")
	})
	router.Compile()

	req, _ = requestAndCheck(t, router, "GET", "/foo/teststr", "get pathtest1")
	if Vars(req)[0] != "teststr" {
		t.Errorf("should set correct var = 'teststr', got: %s", Vars(req)[0])
	}

	req, _ = requestAndCheck(t, router, "GET", "/foo/anothervar", "get pathtest1")
	if Vars(req)[0] != "anothervar" {
		t.Errorf("should set correct var = 'anothervar', got: %s", Vars(req)[0])
	}

	req, _ = requestAndCheck(t, router, "GET", "/foo/pvar1/pvar2", "get pathtest2")
	vars := Vars(req)
	if vars == nil {
		t.Fail()
	}
	if vars[0] != "pvar1" || vars[1] != "pvar2" {
		t.Errorf("should set correct var = [ \"pvar1\", \"pvar2\" ], got: %s", vars)
	}

	req, _ = requestAndCheck(t, router, "GET", "/foo/pvar3/pvar4", "get pathtest2")
	vars = Vars(req)
	if vars == nil {
		t.Fail()
	}
	if vars[0] != "pvar3" || vars[1] != "pvar4" {
		t.Errorf("should set correct var = [ \"pvar3\", \"pvar4\" ], got: %s", vars)
	}
}
