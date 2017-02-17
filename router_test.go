package router

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var DummyHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
})

func TestRegister(t *testing.T) {
	r := New()

	r.Register(http.MethodGet, "one/two/three", DummyHandler)
	r.Register(http.MethodPost, "one/two/three", DummyHandler)

	req := httptest.NewRequest(http.MethodGet, "/one/two/three", nil)
	_, err := r.GetEndpoint(req)
	if err != nil {
		t.Error("Unexpected route error")
	}
}

func TestRegisterRoute(t *testing.T) {
	r := New()

	r.RegisterRoute("one/two/three", Endpoints{
		http.MethodGet:  DummyHandler,
		http.MethodPost: DummyHandler,
	})

	req := httptest.NewRequest(http.MethodGet, "/one/two/three", nil)
	_, err := r.GetEndpoint(req)
	if err != nil {
		t.Error("Unexpected route error")
	}
}

func TestPrint(t *testing.T) {
	r := New()
	r.RegisterRoute("one/two/three", Endpoints{
		http.MethodGet:  DummyHandler,
		http.MethodPost: DummyHandler,
	})
	r.Print()
}

func TestGetEndpointNotFound(t *testing.T) {
	r := New()
	r.RegisterRoute("one/two/three", Endpoints{
		http.MethodGet:  DummyHandler,
		http.MethodPost: DummyHandler,
	})

	req := httptest.NewRequest(http.MethodGet, "/two/three", nil)
	_, err := r.GetEndpoint(req)
	if err == nil || err != ErrNoURLMatch {
		t.Error("expected ErrNoURLMatch from GetEndpoint")
	}
}

func TestGetEndpointNotAllowed(t *testing.T) {
	r := New()
	r.RegisterRoute("one/two/three", Endpoints{
		http.MethodGet:  DummyHandler,
		http.MethodPost: DummyHandler,
	})

	req := httptest.NewRequest(http.MethodPut, "/one/two/three", nil)
	_, err := r.GetEndpoint(req)
	if err == nil || err != ErrNoMethodMatch {
		t.Error("expected ErrNoMethodMatch from GetEndpoint")
	}
}

func TestGetEndpointContext(t *testing.T) {
	r := New()
	r.RegisterRoute("one/{two}/three", Endpoints{
		http.MethodGet:  DummyHandler,
		http.MethodPost: DummyHandler,
	})

	req := httptest.NewRequest(http.MethodGet, "/one/five/three", nil)
	r.GetEndpoint(req)
	ctx := req.Context()
	v := ctx.Value("two")
	if v == nil {
		log.Printf("%+v\n", ctx)
		t.Error("Key not present in context")
	} else if v != "five" {
		t.Error("Incorrect value in context")
	}
}
