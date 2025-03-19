package gee

import (
	"fmt"
	"reflect"
	"testing"
)

func newTestRouter() *router {
	r := newRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/hello/b/c", nil)
	r.addRoute("GET", "/hi/:name", nil)
	r.addRoute("GET", "/assets/*filepath", nil)
	return r
}

func TestParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"p", "*name"})
	if !ok {
		t.Fatal("test parsePattern failed")
	}
}

func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	n, params := r.getRoute("GET", "/hello/geektutu")

	if n == nil {
		t.Fatal("test getRoute failed, nil should not be returned")
	}

	if n.pattern != "/hello/:name" {
		t.Fatal("test getRoute failed, should match /hello/:name")
	}

	if params["name"] != "geektutu" {
		t.Fatal("test getRoute failed, name should be geektutu")
	}

	fmt.Printf("matched path: %s, params['name']: %s\n", n.pattern, params["name"])
}
