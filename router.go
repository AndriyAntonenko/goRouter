package goRouter

import (
	"fmt"
	"net/http"
)

type Handler = func(w http.ResponseWriter, r *http.Request, params *RouterParams)

type Router struct {
	MethodToHandlers map[string]*RouterTrie
}

func NewRouter() *Router {
	return &Router{
		MethodToHandlers: map[string]*RouterTrie{},
	}
}

func (r *Router) Post(path string, handler Handler) *Router {
	return r.addHandler(http.MethodPost, path, handler)
}

func (r *Router) Get(path string, handler Handler) *Router {
	return r.addHandler(http.MethodGet, path, handler)
}

func (r *Router) Put(path string, handler Handler) *Router {
	return r.addHandler(http.MethodPut, path, handler)
}

func (r *Router) Patch(path string, handler Handler) *Router {
	return r.addHandler(http.MethodPatch, path, handler)
}

func (r *Router) Delete(path string, handler Handler) *Router {
	return r.addHandler(http.MethodDelete, path, handler)
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	trie := router.getMethodTrie(r.Method)
	ps := NewRouterParams()
	node := trie.Lookup(r.URL.Path, ps)

	if node == nil || node.handler == nil {
		router.defaultHandler(w, r)
		return
	}

	node.handleCall(w, r, ps)
}

func (router *Router) defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(fmt.Sprintf("404 Not Found; %s %s", r.Method, r.URL.Path)))
}

func (r *Router) getMethodTrie(method string) *RouterTrie {
	trie, ok := r.MethodToHandlers[method]
	if !ok {
		trie = NewRouterTrie("*", 10, "")
		r.MethodToHandlers[method] = trie
	}

	return trie
}

func (r *Router) addHandler(method string, path string, handler Handler) *Router {
	trie := r.getMethodTrie(method)
	trie.AddNode(path, handler)
	return r
}
