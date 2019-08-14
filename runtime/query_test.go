package runtime

import (
	"strings"
	"testing"
)

func TestToQueryURL(t *testing.T) {
	type Foo struct {
		A string `json:"a"`
		B string `json:"b"`
	}

	foo := &Foo{A: "a", B: "b"}
	actual := toQueryURL(DefaultQueryEncoder, "http://example.com?b=1&c=c", foo)

	if !strings.Contains(actual, "a=a") {
		t.Fatalf("want a, got %s", actual)
	}
	if !strings.Contains(actual, "b=b") {
		t.Fatalf("want b, got %s", actual)
	}
	if !strings.Contains(actual, "c=c") {
		t.Fatalf("want c, got %s", actual)
	}
}
