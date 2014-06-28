// github.com/c9s/r3 binding for go
package r3

/*
#cgo LDFLAGS: -lr3

#include <r3/r3.h>
#include <stdlib.h>

void * getData(route *r) {
	return r->data;
}

int getVarsLength(match_entry *e) {
	return e->vars->len;
}

int getRequestMethod(match_entry *e) {
	return e->request_method;
}

const char* getVar(match_entry *e, int i) {
	return e->vars->tokens[i];
}
*/
import "C"
import "unsafe"
import (
	"errors"
)

var (
	MethodGet     Method = C.METHOD_GET
	MethodPost    Method = C.METHOD_POST
	MethodPut     Method = C.METHOD_PUT
	MethodDelete  Method = C.METHOD_DELETE
	MethodPatch   Method = C.METHOD_PATCH
	MethodHead    Method = C.METHOD_HEAD
	MethodOptions Method = C.METHOD_OPTIONS
)

type Method int
type Node struct {
	node *C.node
}
type Data struct {
	Value interface{}
}

type Tree Node
type Router Node

// Create a new Tree
func NewTree(capacity int) *Tree {
	var n *C.node
	n = C.r3_tree_create(C.int(capacity))
	t := &Tree{n}
	return t
}

// Insert route with method and arbitary data
func (n *Tree) InsertRoute(m Method, path string, data interface{}) {
	d := &Data{
		Value: data,
	}
	p := C.CString(path)
	C.r3_tree_insert_routel(n.node, C.int(m), p, C.int(C.strlen(p)), unsafe.Pointer(d))
}

// Insert path
func (n *Tree) InsertPath(path string, data interface{}) {
	d := &Data{
		Value: data,
	}
	p := C.CString(path)
	C.r3_tree_insert_pathl(n.node, p, C.int(C.strlen(p)), unsafe.Pointer(d))
}

// Compile the tree, returns error if any
func (n *Tree) Compile() error {
	var err C.int
	var errstr *C.char
	err = C.r3_tree_compile(n.node, &errstr)
	if err != 0 {
		defer C.free(unsafe.Pointer(errstr))
		return errors.New(C.GoString(errstr))
	}
	return nil
}

// Match a MatchEntry, returns a Route object if found, nil otherwise
func (n *Tree) MatchRoute(e *MatchEntry) *Route {
	r := C.r3_tree_match_route(n.node, e.entry)
	if r == nil {
		return nil
	}
	return &Route{
		route: r,
	}
}

// Free the memory of a tree
func (n *Tree) Free() {
	C.r3_tree_free(n.node)
	n.node = nil
}

// Dump tree to stdout
func (n *Tree) Dump() {
	C.r3_tree_dump(n.node, 0)
}

type MatchEntry struct {
	entry *C.match_entry
}

// creates a new MatchEntry
func NewMatchEntry(path string) *MatchEntry {
	p := C.CString(path)
	e := C.match_entry_createl(p, C.int(C.strlen(p)))
	return &MatchEntry{
		entry: e,
	}
}

func (e *MatchEntry) RequestMethod() Method {
	return Method(C.getRequestMethod(e.entry))
}

// Set request method of the entry
func (e *MatchEntry) SetRequestMethod(m Method) {
	e.entry.request_method = C.int(m)
}

// Get tokens in the path
func (e *MatchEntry) Vars() *[]string {
	var length int = int(C.getVarsLength(e.entry))
	var tokens []string = make([]string, length)
	for i := 0; i < length; i++ {
		tokens[i] = C.GoString(C.getVar(e.entry, C.int(i)))
	}
	return &tokens
}

// Free memory
func (e *MatchEntry) Free() {
	C.match_entry_free(e.entry)
	e.entry = nil
}

type Route struct {
	route *C.route
}

// Returns data payload of the route
func (r *Route) Data() interface{} {
	var p *Data
	p = (*Data)(C.getData(r.route))
	return p.Value
}

// Free memory
func (r *Route) Free() {
	C.r3_route_free(r.route)
	r.route = nil
}
