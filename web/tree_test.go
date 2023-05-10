package web

import (
	"reflect"
	"testing"
)

func newNode() *node {
	return &node{
		children: make([]*node, 0),
	}
}

func TestParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	if !ok {
		t.Fatal("test parsePattern failed")
	}
}

func TestGetRoute(t *testing.T) {
	n := newNode()

	n.insert("/p/:name", parsePattern("/p/:name"), 0, nil)
	n.insert("/p/a", parsePattern("/p/a"), 0, nil)
	n.insert("/p/a/*", parsePattern("/p/a/*"), 0, nil)
	n.insert("/p/b", parsePattern("/p/b"), 0, nil)
	n.insert("/p/b/c", parsePattern("/p/b/c"), 0, nil)
	n.insert("/p/b/:name/:age", parsePattern("/p/b/:name/:age"), 0, nil)

	res := n.search(parsePattern("/p/a"), 0)
	if res.pattern != "/p/a" {
		t.Fatal("should match /p/a, but got", res.pattern)
	}

	res = n.search(parsePattern("/p/a/b/c"), 0)
	if res.pattern != "/p/a/*" {
		t.Fatal("should match /p/a/*, but got", res.pattern)
	}

	res = n.search(parsePattern("/p/b/c"), 0)
	if res.pattern != "/p/b/c" {
		t.Fatal("should match /p/b/c, but got", res.pattern)
	}

	res = n.search(parsePattern("/p/d"), 0)
	if res.pattern != "/p/:name" {
		t.Fatal("should match /p/:name, but got", res.pattern)
	}

	res = n.search(parsePattern("/p/b/bob/18"), 0)
	if res.pattern != "/p/b/:name/:age" {
		t.Fatal("should match /p/b/:name/:age, but got", res.pattern)
	}

	// n.insert("/p/*", parsePattern("/p/d/*"), 0, nil)
	// n.insert("/p/*name/*", parsePattern("/p/*name/*"), 0, nil)

}
