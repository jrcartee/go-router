package router

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var ErrNoURLMatch = errors.New("No matching URL found in routes")
var ErrNoMethodMatch = errors.New("No matching Method found in endpoints")

type Router struct {
	root *RouteNode
}

func New() *Router {
	root := NewRouteNode("")
	return &Router{root}
}

func (self *Router) RegisterRoute(p string, e Endpoints) {
	p_arr := strings.Split(p, "/")
	p_arr = removeEmpty(p_arr)

	current := self.root
	for _, v := range p_arr {
		current = current.GetOrCreate(v)
	}
	current.InsertEndpoints(e)
}

func (self *Router) Register(method string, path string, h http.Handler) {
	p_arr := strings.Split(path, "/")
	p_arr = removeEmpty(p_arr)

	current := self.root
	for _, v := range p_arr {
		current = current.GetOrCreate(v)
	}
	current.InsertEndpoints(Endpoints{
		method: h,
	})
}

func (self *Router) Print() {
	fmt.Println("Routing Tree:")
	printTree(self.root, "")
}

func (self *Router) GetEndpoint(r *http.Request) (http.Handler, error) {
	p_arr := strings.Split(r.URL.String(), "/")
	p_arr = removeEmpty(p_arr)
	ctx := make(map[string]string, 0)

	// iterate path segments and traverse tree for matches
	current := self.root
	for _, p := range p_arr {
		current = current.FindMatch(p)
		if current == nil {
			return nil, ErrNoURLMatch
		}
		if current.matcher != nil {
			// this node is a param'd segment, store context
			k, v := current.matcher.key, p
			ctx[k] = v
		}
	}
	// check for http Method
	h, ok := (*current.endpoints)[r.Method]
	if !ok {
		return nil, ErrNoMethodMatch
	}
	// add context to request
	rc := r.Context()
	for k, v := range ctx {
		rc = context.WithValue(rc, k, v)
	}
	*r = *r.WithContext(rc)
	return h, nil
}

func removeEmpty(s []string) []string {
	var r []string
	for _, v := range s {
		if v != "" {
			r = append(r, v)
		}
	}
	return r
}

func printTree(node *RouteNode, indent string) {

	endpoints := *node.endpoints

	var allowed_methods string
	if len(endpoints) != 0 {
		keys := make([]string, len(endpoints))
		i := 0
		for k := range endpoints {
			keys[i] = k
			i++
		}

		key_str := strings.Join(keys, ", ")
		allowed_methods = fmt.Sprintf("[%s]", key_str)
	} else {
		allowed_methods = ""
	}

	if node.Path != "" {
		fmt.Printf("%s └──%s\t%s\n", indent, node.Path, allowed_methods)
		indent = "   " + indent
	}

	for _, n := range node.children {
		printTree(n, indent)
	}
}
