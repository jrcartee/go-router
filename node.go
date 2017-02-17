package router

import (
	"net/http"
	"regexp"
	// "fmt"
)

var url_param_matcher *regexp.Regexp
var url_param_regex_matcher *regexp.Regexp

func init() {
	url_param_matcher = regexp.MustCompile(`{(.+)}`)
	url_param_regex_matcher = regexp.MustCompile(`{([A-Za-z]+):(.+)}`)
}

const (
	GenericParam = iota
	RegexpParam
)

type ParamMatcher struct {
	mtype int
	key   string
	exp   *regexp.Regexp
}

func NewParamMatcher(path string) *ParamMatcher {

	if url_param_regex_matcher.MatchString(path) {
		submatch := url_param_regex_matcher.FindStringSubmatch(path)
		exp, err := regexp.Compile(submatch[2])
		if err != nil {
			panic("Invalid regexp supplied in endpoint path: " + path)
		}
		return &ParamMatcher{RegexpParam, submatch[1], exp}

	} else if url_param_matcher.MatchString(path) {
		key := url_param_matcher.FindStringSubmatch(path)[1]
		return &ParamMatcher{GenericParam, key, nil}

	} else {
		return nil
	}
}

type Endpoints map[string]http.Handler
type RouteNode struct {
	Path      string
	matcher   *ParamMatcher
	endpoints *Endpoints
	children  []*RouteNode
}

func NewRouteNode(p string) *RouteNode {
	m := NewParamMatcher(p)
	e := make(Endpoints, 0)
	c := make([]*RouteNode, 0)
	return &RouteNode{p, m, &e, c}
}

func (n *RouteNode) GetOrCreate(p string) (r *RouteNode) {
	for _, v := range n.children {
		if v.Path == p {
			return v
		}
	}
	tmp := NewRouteNode(p)
	n.children = append(n.children, tmp)
	return tmp
}

func (n *RouteNode) InsertEndpoints(e Endpoints) {
	tmp := (*n.endpoints)
	for k, v := range e {
		tmp[k] = v
	}
	n.endpoints = &tmp
}

func (n *RouteNode) MatchSubPath(p string) bool {
	switch {
	case n.matcher == nil && n.Path == p:
		return true

	case n.matcher != nil && n.matcher.mtype == GenericParam:
		return true

	case n.matcher != nil && n.matcher.mtype == RegexpParam && n.matcher.exp.MatchString(p):
		return true
	}
	return false
}

func (n *RouteNode) FindMatch(p string) *RouteNode {
	var best_match *RouteNode
	for _, v := range n.children {
		match := v.MatchSubPath(p)
		if !match {
			continue
		}
		if v.matcher == nil {
			return v // literal match, highest priority
		}
		if v.matcher.mtype == RegexpParam {
			best_match = v
			continue // regexp match, second priority
		}
		if best_match == nil {
			best_match = v // generic match, lowest priority
		}
	}
	return best_match
}

func (n *RouteNode) GetContext(p string) (key, value string) {
	if n.matcher == nil {
		panic("GetContext called on RouteNode with no matcher")
	}

	return n.matcher.key, p

}
