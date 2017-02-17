package router

import (
	"testing"
)

func TestCreate(t *testing.T) {
	node := NewRouteNode("")
	node.GetOrCreate("inner")
	match := node.FindMatch("inner")
	if match == nil {
		t.Error("child node not found")
	}
}

func TestGet(t *testing.T) {
	node := NewRouteNode("")
	node.GetOrCreate("inner")
	node.GetOrCreate("inner")
	node.GetOrCreate("inner")
	if len(node.children) != 1 {
		t.Error("duplicate routes created")
	}
}

func TestMatchSubPath(t *testing.T) {
	node := NewRouteNode("")
	n1 := node.GetOrCreate("{id}")
	n2 := node.GetOrCreate("inner")
	n3 := n2.GetOrCreate(`{id:\d+}`)

	var match bool

	match = node.MatchSubPath("")
	if !match {
		t.Error("match expected")
	}
	match = n1.MatchSubPath("anything")
	if !match {
		t.Error("match expected")
	}
	match = n2.MatchSubPath("inner")
	if !match {
		t.Error("match expected")
	}
	match = n3.MatchSubPath("123")
	if !match {
		t.Error("match expected")
	}
}

func TestFindMatch(t *testing.T) {
	node := NewRouteNode("")
	na := node.GetOrCreate("{id}")
	nb := node.GetOrCreate(`{id:\d+}`)
	nc := node.GetOrCreate("inner")

	var match *RouteNode

	match = node.FindMatch("inner")
	if match != nc {
		t.Error("Literal not prioritized")
	}

	match = node.FindMatch("123")
	if match != nb {
		t.Error("Regexp not prioritized")
	}

	match = node.FindMatch("else")
	if match != na {
		t.Error("unexpected match")
	}
}

func TestBadRegex(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Error("Expected Panic")
		}
	}()

	NewParamMatcher("{key:][}")
}

func TestUnexpectedGetContext(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Error("Expected Panic")
		}
	}()

	node := NewRouteNode("")
	node.GetContext("example")
}

func TestExpectedGetContext(t *testing.T) {
	node := NewRouteNode("{id}")
	k, v := node.GetContext("whatever")

	if k != "id" || v != "whatever" {
		t.Error("Incorrect key/value returned from GetContext")
	}

	node = NewRouteNode("{key:[a-z]+}")
	k, v = node.GetContext("regex")

	if k != "key" || v != "regex" {
		t.Error("Incorrect key/value returned from GetContext")
	}
}
