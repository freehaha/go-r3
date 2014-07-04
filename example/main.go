package main

import (
	"fmt"
	"github.com/freehaha/go-r3"
	"github.com/freehaha/go-r3/mux"
	"log"
	"net/http"
)

func main() {
	/* instanciate router */
	r := mux.NewRouter()
	/* optional, it should be GCed automatically */
	defer r.Free()

	/* static paths */
	r.HandleFunc(r3.MethodGet, "/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "hello world")
	})

	r.HandleFunc(r3.MethodGet, "/foo/bar", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "foo/bar")
	})

	/* path with parameters */
	r.HandleFunc(r3.MethodGet, "/path/{id}", func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		fmt.Fprintf(w, "path, args: %s", vars)
	})

	r.HandleFunc(r3.MethodGet, "/path/{id}/{arg2}", func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		fmt.Fprintf(w, "path, args: %s", vars)
	})

	/* must be compiled before use */
	r.Compile()

	log.Println("listening")
	err := http.ListenAndServe(":3003", r)
	if err != nil {
		log.Printf("%s", err.Error())
	}
}
