// Copyright 2013 Julien Schmidt. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

// Package httprouter is a trie based high performance HTTP request router.
//
// A trivial example is:
//
//	package main
//
//	import (
//	    "fmt"
//	    "github.com/julienschmidt/httprouter"
//	    "net/http"
//	    "log"
//	)
//
//	func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
//	    fmt.Fprint(w, "Welcome!\n")
//	}
//
//	func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	    fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
//	}
//
//	func main() {
//	    router := httprouter.New()
//	    router.GET("/", Index)
//	    router.GET("/hello/:name", Hello)
//
//	    log.Fatal(http.ListenAndServe(":8080", router))
//	}
//
// The router matches incoming requests by the request method and the path.
// If a handle is registered for this path and method, the router delegates the
// request to that function.
// For the methods GET, POST, PUT, PATCH and DELETE shortcut functions exist to
// register handles, for all other methods router.ServiceInfo can be used.
//
// The registered path, against which the router matches incoming requests, can
// contain two types of parameters:
//
//	Syntax    Type
//	:name     named parameter
//	*name     catch-all parameter
//
// Named parameters are dynamic path segments. They match anything until the
// next '/' or the path end:
//
//	Path: /blog/:category/:post
//
//	Requests:
//	 /blog/go/request-routers            match: category="go", post="request-routers"
//	 /blog/go/request-routers/           no match, but the router would redirect
//	 /blog/go/                           no match
//	 /blog/go/request-routers/comments   no match
//
// Catch-all parameters match anything until the path end, including the
// directory index (the '/' before the catch-all). Since they match anything
// until the end, catch-all parameters must always be the final path element.
//
//	Path: /files/*filepath
//
//	Requests:
//	 /files/                             match: filepath="/"
//	 /files/LICENSE                      match: filepath="/LICENSE"
//	 /files/templates/article.html       match: filepath="/templates/article.html"
//	 /files                              no match, but the router would redirect
//
// The value of parameters is saved as a slice of the Param struct, consisting
// each of a key and a value. The slice is passed to the ServiceInfo func as a third
// parameter.
// There are two ways to retrieve the value of a parameter:
//
//	// by the name of the parameter
//	user := ps.ByName("user") // defined by :user or *user
//
//	// by the index of the parameter. This way you can also get the name (key)
//	thirdKey   := ps[2].Key   // the name of the 3rd parameter
//	thirdValue := ps[2].Value // the value of the 3rd parameter
package httprouter

import (
	v1 "github.com/baker-yuan/go-gateway/go-gateway-admin/api/admin/v1"
	"net/http"
	"strings"
)

// ServiceInfo is a function that can be registered to a route to handle HTTP
// requests. Like http.HandlerFunc, but has a third parameter for the values of
// wildcards (variables).
type ServiceInfo struct {
	ID            uint32           // ID
	GwUrl         string           // 网关接口
	HttpType      string           // 接口类型 net/http/method.go
	Status        v1.Status        // 接口状态
	Application   string           // 应用名称
	InterfaceType v1.InterfaceType // 接口协议
	Config        string           // 指定协议的配置
	InterfaceUrl  string           // 接口方法
}

// Param is a single URL parameter, consisting of a key and a value.
type Param struct {
	Key   string
	Value string
}

// Params is a Param-slice, as returned by the router.
// The slice is ordered, the first URL parameter is also the first slice value.
// It is therefore safe to read values by the index.
type Params []Param

// ByName returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (ps Params) ByName(name string) string {
	for i := range ps {
		if ps[i].Key == name {
			return ps[i].Value
		}
	}
	return ""
}

// RouterRepo is a http.Handler which can be used to dispatch requests to different
// handler functions via configurable routes
type RouterRepo struct {
	trees map[string]*node

	// Cached value of global (*) allowed methods
	globalAllowed string
}

// New returns a new initialized RouterRepo.
// Path auto-correction, including trailing slashes, is enabled by default.
func New() *RouterRepo {
	return &RouterRepo{}
}

// GET is a shortcut for router.ServiceInfo(http.MethodGet, path, handle)
func (r *RouterRepo) GET(path string, handle *ServiceInfo) {
	r.Handle(http.MethodGet, path, handle)
}

// HEAD is a shortcut for router.ServiceInfo(http.MethodHead, path, handle)
func (r *RouterRepo) HEAD(path string, handle *ServiceInfo) {
	r.Handle(http.MethodHead, path, handle)
}

// OPTIONS is a shortcut for router.ServiceInfo(http.MethodOptions, path, handle)
func (r *RouterRepo) OPTIONS(path string, handle *ServiceInfo) {
	r.Handle(http.MethodOptions, path, handle)
}

// POST is a shortcut for router.ServiceInfo(http.MethodPost, path, handle)
func (r *RouterRepo) POST(path string, handle *ServiceInfo) {
	r.Handle(http.MethodPost, path, handle)
}

// PUT is a shortcut for router.ServiceInfo(http.MethodPut, path, handle)
func (r *RouterRepo) PUT(path string, handle *ServiceInfo) {
	r.Handle(http.MethodPut, path, handle)
}

// PATCH is a shortcut for router.ServiceInfo(http.MethodPatch, path, handle)
func (r *RouterRepo) PATCH(path string, handle *ServiceInfo) {
	r.Handle(http.MethodPatch, path, handle)
}

// DELETE is a shortcut for router.ServiceInfo(http.MethodDelete, path, handle)
func (r *RouterRepo) DELETE(path string, handle *ServiceInfo) {
	r.Handle(http.MethodDelete, path, handle)
}

// Handle
// ServiceInfo registers a new request handle with the given path and method.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (r *RouterRepo) Handle(method, path string, handle *ServiceInfo) {
	if len(path) < 1 || path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}

	if r.trees == nil {
		r.trees = make(map[string]*node)
	}

	root := r.trees[method]
	if root == nil {
		root = new(node)
		r.trees[method] = root

		r.globalAllowed = r.allowed("*", "")
	}

	root.addRoute(path, handle)
}

// Lookup allows the manual lookup of a method + path combo.
// This is e.g. useful to build a framework around this router.
// If the path was found, it returns the handle function and the path parameter
// values. Otherwise the third return value indicates whether a redirection to
// the same path with an extra / without the trailing slash should be performed.
func (r *RouterRepo) Lookup(method, path string) (*ServiceInfo, Params, bool) {
	if root := r.trees[method]; root != nil {
		return root.getValue(path)
	}
	return nil, nil, false
}

func (r *RouterRepo) allowed(path, reqMethod string) (allow string) {
	allowed := make([]string, 0, 9)

	if path == "*" { // server-wide
		// empty method is used for internal calls to refresh the cache
		if reqMethod == "" {
			for method := range r.trees {
				if method == http.MethodOptions {
					continue
				}
				// Add request method to list of allowed methods
				allowed = append(allowed, method)
			}
		} else {
			return r.globalAllowed
		}
	} else { // specific path
		for method := range r.trees {
			// Skip the requested method - we already tried this one
			if method == reqMethod || method == http.MethodOptions {
				continue
			}

			handle, _, _ := r.trees[method].getValue(path)
			if handle != nil {
				// Add request method to list of allowed methods
				allowed = append(allowed, method)
			}
		}
	}

	if len(allowed) > 0 {
		// Add request method to list of allowed methods
		allowed = append(allowed, http.MethodOptions)

		// Sort allowed methods.
		// sort.Strings(allowed) unfortunately causes unnecessary allocations
		// due to allowed being moved to the heap and interface conversion
		for i, l := 1, len(allowed); i < l; i++ {
			for j := i; j > 0 && allowed[j] < allowed[j-1]; j-- {
				allowed[j], allowed[j-1] = allowed[j-1], allowed[j]
			}
		}

		// return as comma separated list
		return strings.Join(allowed, ", ")
	}
	return
}
