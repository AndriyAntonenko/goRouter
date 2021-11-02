package goRouter

import (
	"fmt"
	"net/http"
)

type Handler = func(w http.ResponseWriter, r *http.Request, params *RouterParams)

type RouterCors struct {
	Origins string
	Methods string
	Headers string
	MaxAge  string
}

type Router struct {
	methodToHandlers map[string]*RouterTrie
	corsConf         *RouterCors
}

func NewRouter() *Router {
	return &Router{
		methodToHandlers: map[string]*RouterTrie{},
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

func (r *Router) EnableCors(cors *RouterCors) {
	r.corsConf = cors
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if router.corsConf != nil && r.Method == http.MethodOptions {
		router.corsHandler(w, r)
		return
	}

	trie := router.getMethodTrie(r.Method)
	ps := NewRouterParams()
	node := trie.Lookup(r.URL.Path, ps)

	if node == nil || node.handler == nil {
		router.defaultHandler(w, r)
		return
	}

	if router.corsConf != nil {
		w.Header().Set("Access-Control-Allow-Origin", router.corsConf.Origins)
	}
	node.handleCall(w, r, ps)
}

func (router *Router) defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(fmt.Sprintf("404 Not Found; %s %s", r.Method, r.URL.Path)))
}

func (r *Router) getMethodTrie(method string) *RouterTrie {
	trie, ok := r.methodToHandlers[method]
	if !ok {
		trie = NewRouterTrie("*", 10, "")
		r.methodToHandlers[method] = trie
	}

	return trie
}

// Handler for preflights requests
func (router *Router) corsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", router.corsConf.Origins)
	w.Header().Set("Access-Control-Allow-Methods", router.corsConf.Methods)
	w.Header().Set("Access-Control-Allow-Headers", router.corsConf.Headers)
	w.Header().Set("Access-Control-Max-Age", router.corsConf.MaxAge)
	w.WriteHeader(http.StatusNoContent)
}

func (r *Router) addHandler(method string, path string, handler Handler) *Router {
	trie := r.getMethodTrie(method)
	trie.AddNode(path, handler)
	return r
}
