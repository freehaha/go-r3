package router

import (
	"github.com/freehaha/go-r3"
	"github.com/gorilla/context"
	"log"
	"net/http"
	"runtime"
)

type Router struct {
	tree            *r3.Tree
	NotFoundHandler http.Handler
}

var methods map[string]r3.Method = map[string]r3.Method{
	"GET":     r3.MethodGet,
	"PUT":     r3.MethodPut,
	"POST":    r3.MethodPost,
	"DELETE":  r3.MethodDelete,
	"OPTIONS": r3.MethodOptions,
	"PATCH":   r3.MethodPatch,
}

/* instanciate a new Router */
func NewRouter() *Router {
	r := &Router{}
	r.tree = r3.NewTree(10)
	runtime.SetFinalizer(r, finalizeRouter)
	return r
}

func finalizeRouter(r *Router) {
	log.Print("finalizing router")
	r.tree.Free()
}

func (r *Router) Compile() error {
	return r.tree.Compile()
}

func (r *Router) Free() {
	finalizeRouter(r)
}

/* insert a path to the router with specific handler */
func (r *Router) HandleFunc(method r3.Method, path string, handler http.HandlerFunc) {
	r.tree.InsertRoute(method, path, handler)
}

/* implement Mux */
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ent := r3.NewMatchEntry(req.URL.Path)
	ent.SetRequestMethod(methods[req.Method])
	log.Print(req.URL.Path)
	log.Print(req.Method)
	defer ent.Free()
	route := r.tree.MatchRoute(ent)
	if route != nil {
		context.Set(req, "vars", ent.Vars())
		handler := route.Data().(http.HandlerFunc)
		handler(w, req)
	} else {
		handler := r.NotFoundHandler
		if handler == nil {
			handler = http.NotFoundHandler()
		}
		handler.ServeHTTP(w, req)
	}
}
