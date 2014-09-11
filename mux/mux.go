package mux

import (
	"github.com/freehaha/go-r3"
	"github.com/gorilla/context"
	"net/http"
	"runtime"
)

type Router struct {
	tree            *r3.Tree
	NotFoundHandler http.Handler
}

type NegroniRouter struct {
	Router
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
	r.tree.Free()
}

func finalizeNRouter(r *NegroniRouter) {
	r.tree.Free()
}

func (r *Router) Compile() error {
	return r.tree.Compile()
}

func (r *Router) Free() {
	finalizeRouter(r)
}

/* Helper function for HandleFunc(r3.MethodGet, path, handler) */
func (r *Router) Get(path string, handler http.HandlerFunc) {
	r.HandleFunc(r3.MethodGet, path, handler)
}

/* Helper function for HandleFunc(r3.MethodPost, path, handler) */
func (r *Router) Post(path string, handler http.HandlerFunc) {
	r.HandleFunc(r3.MethodPost, path, handler)
}

/* Helper function for HandleFunc(r3.MethodPut, path, handler) */
func (r *Router) Put(path string, handler http.HandlerFunc) {
	r.HandleFunc(r3.MethodPut, path, handler)
}

/* Helper function for HandleFunc(r3.MethodPatch, path, handler) */
func (r *Router) Patch(path string, handler http.HandlerFunc) {
	r.HandleFunc(r3.MethodPatch, path, handler)
}

/* Helper function for HandleFunc(r3.MethodDelete, path, handler) */
func (r *Router) Delete(path string, handler http.HandlerFunc) {
	r.HandleFunc(r3.MethodDelete, path, handler)
}

/* Helper function for HandleFunc(r3.MethodOptions, path, handler) */
func (r *Router) Options(path string, handler http.HandlerFunc) {
	r.HandleFunc(r3.MethodOptions, path, handler)
}

/* returns the variables captured by the placeholders. */
func Vars(r *http.Request) []string {
	data := context.Get(r, "vars")
	if data == nil {
		return nil
	}
	return context.Get(r, "vars").([]string)
}

/* insert a path to the router with specific handlerfunc */
func (r *Router) HandleFunc(method r3.Method, path string, handler http.HandlerFunc) {
	r.Handle(method, path, http.HandlerFunc(handler))
}

/* insert a path to the router with specific handler */
func (r *Router) Handle(method r3.Method, path string, handler http.Handler) {
	r.tree.InsertRoute(method, path, handler)
	/* keep a reference on it so it doens't get GCed */
}

/* implement Mux */
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ent := r3.NewMatchEntry(req.URL.Path)
	ent.SetRequestMethod(methods[req.Method])
	defer ent.Free()
	route := r.tree.MatchRoute(ent)
	if route != nil {
		context.Set(req, "vars", ent.Vars())
		handler := route.Data().(http.Handler)
		handler.ServeHTTP(w, req)
		context.Clear(req)
	} else {
		handler := r.NotFoundHandler
		if handler == nil {
			handler = http.NotFoundHandler()
		}
		handler.ServeHTTP(w, req)
	}
}

func NewNegroniRouter() *NegroniRouter {
	r := &NegroniRouter{
		Router{},
	}
	r.tree = r3.NewTree(10)
	runtime.SetFinalizer(r, finalizeNRouter)
	return r
}

/* implement negroni middleware */
func (r *NegroniRouter) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	ent := r3.NewMatchEntry(req.URL.Path)
	ent.SetRequestMethod(methods[req.Method])
	defer ent.Free()
	route := r.tree.MatchRoute(ent)
	if route != nil {
		context.Set(req, "vars", ent.Vars())
		handler := route.Data().(http.Handler)
		handler.ServeHTTP(w, req)
		context.Clear(req)
	} else {
		/* skip to next middleware */
		next(w, req)
	}
}
