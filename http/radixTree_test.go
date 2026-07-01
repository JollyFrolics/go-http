package http

import (
	"codeberg.org/JollyFrolics/go-http.git/colorize"
	"testing"
)

func assertPanics(t *testing.T, f func(), msg string) {
	t.Helper()

	defer func() {
		if r := recover(); r == nil {
			t.Error(msg)
		}
	}()

	f()
}

func TestRadixNode_insert(t *testing.T) {
	tree := &RadixNode{}
	assertPanics(
		t, func() {
			_ = tree.insert("", 100)

		}, "Inserting empty string should panic didn't",
	)

	curr := tree.insert("/", 100)

	if curr == nil {
		t.Error(colorize.Error("Expected a tree node got"), curr)
	}

	if curr.Data != 100 {
		t.Error(colorize.Error("Expected data of 100 got"), curr.Data)
	}

	if tree != curr {
		t.Error(colorize.Error("Expected the map of the \"/\" to be the one returned from insert call wasn't"))
	}

	if tree.Data != 100 {
		t.Error(colorize.Error("Expected data of 100 got"), curr.Data)
	}

	curr = tree.insert("/home", 200)

	if curr.Data != 200 {
		t.Error(colorize.Error("Expected data of 200 got"), curr.Data)
	}

	if tree.Data != 100 {
		t.Error(colorize.Error("Expected data of 100 got"), curr.Data)
	}

	if tree.Child['h'] != curr {
		t.Error(colorize.Error("Expected maping to result in home got: "), curr)
	}

	curr = tree.insert("/home", 100)

	if curr.Data != 100 {
		t.Error(colorize.Error("Expected data of 100 got"), curr.Data)
	}

	if tree.Data != 100 {
		t.Error(colorize.Error("Expected data of 100 got"), curr.Data)
	}

	if tree.Child['h'] != curr {
		t.Error(colorize.Error("Expected maping to result in /home got: "), curr)
	}

	curr = tree.insert("/home", 200)
	curr = tree.insert("/high", 300)

	if curr.Data != 300 {
		t.Error(colorize.Error("Expected data of 300 got"), curr.Data)
	}

	if tree.Data != 100 {
		t.Error(colorize.Error("Expected data of 100 got"), curr.Data)
	}

	if tree.Child['h'].Child['o'].Data != 200 {
		t.Error(colorize.Error("Expected data of 200 got"), curr.Data)
	}

	if tree.Child['h'].Child['i'].Data != 300 {
		t.Error(colorize.Error("Expected data of 300 got"), curr.Data)
	}

	if tree.Child['h'].Child['i'] != curr {
		t.Error(colorize.Error("Mapping of /, h, i is incorrect"), curr.Data)
	}

	root := &RadixNode{}

	root.insert("abcd", "first")

	root.insert("ab", "second")
}

func Test_longestCommonPrefix(t *testing.T) {
	tests := []struct {
		s    string
		p    string
		want int
	}{
		{"abcd", "abcd", 4},
		{"abcd", "abef", 2},
		{"abcd", "ab", 2},
		{"ab", "abcd", 2},
		{"abcd", "xyz", 0},
		{"", "", 0},
		{"abc", "", 0},
		{"", "abc", 0},
	}

	for _, tt := range tests {
		got := longestCommonPrefix(tt.s, tt.p)

		if got != tt.want {
			t.Errorf(
				"longestCommonPrefix(%q,%q)=%d want %d",
				tt.s,
				tt.p,
				got,
				tt.want,
			)
		}
	}
}

func TestRadixNode_lookup(t *testing.T) {
	tree := &RadixNode{}
	tree.addRoute("/home/:user/*files", 100)

	ctx, err := tree.lookup("/home/protocol/http/radix.go")

	if err != nil {
		t.Fatal("Error occured in lookup" + err.Error())
	}

	if ctx.Params["user"] != "protocol" {
		t.Error("Expected user to be protocol got", ctx.Params["user"])
	}

	if ctx.Params["files"] != "/http/radix.go" {
		t.Error("Expected files to be protocol got", ctx.Params["files"])
	}

	if ctx.Data != 100 {
		t.Error("Expected data to be 100 got", ctx.Data)
	}
}

func TestRadixNode_lookup_ComplexRoute(t *testing.T) {
	tree := &RadixNode{}

	tree.addRoute(
		"/api/:version/users/:userID/projects/:projectID/files/*path",
		200,
	)

	ctx, err := tree.lookup(
		"/api/v1/users/123/projects/456/files/src/router/radix.go",
	)

	if err != nil {
		t.Fatal("Error occurred in lookup: " + err.Error())
	}

	if ctx.Params["version"] != "v1" {
		t.Errorf(
			"Expected version to be v1, got %q",
			ctx.Params["version"],
		)
	}

	if ctx.Params["userID"] != "123" {
		t.Errorf(
			"Expected userID to be 123, got %q",
			ctx.Params["userID"],
		)
	}

	if ctx.Params["projectID"] != "456" {
		t.Errorf(
			"Expected projectID to be 456, got %q",
			ctx.Params["projectID"],
		)
	}

	if ctx.Params["path"] != "/src/router/radix.go" {
		t.Errorf(
			"Expected path to be /src/router/radix.go, got %q",
			ctx.Params["path"],
		)
	}

	if ctx.Data != 200 {
		t.Errorf(
			"Expected data to be 200, got %v",
			ctx.Data,
		)
	}
}

func TestRadixNode_lookup_MultipleRoutes(t *testing.T) {
	tree := &RadixNode{}

	tree.addRoute("/", 1)
	tree.addRoute("/about", 2)
	tree.addRoute("/users/:id", 3)
	tree.addRoute("/users/:id/posts/:postID", 4)
	tree.addRoute("/assets/*file", 5)

	tests := []struct {
		path   string
		data   int
		params map[string]string
	}{
		{
			path:   "/",
			data:   1,
			params: map[string]string{},
		},
		{
			path:   "/about",
			data:   2,
			params: map[string]string{},
		},
		{
			path: "/users/42",
			data: 3,
			params: map[string]string{
				"id": "42",
			},
		},
		{
			path: "/users/42/posts/99",
			data: 4,
			params: map[string]string{
				"id":     "42",
				"postID": "99",
			},
		},
		{
			path: "/assets/css/main.css",
			data: 5,
			params: map[string]string{
				"file": "/css/main.css",
			},
		},
	}

	for _, tt := range tests {
		ctx, err := tree.lookup(tt.path)
		if err != nil {
			t.Fatalf("lookup(%q) failed: %v", tt.path, err)
		}

		if ctx.Data != tt.data {
			t.Errorf(
				"lookup(%q): expected data %d, got %v",
				tt.path,
				tt.data,
				ctx.Data,
			)
		}

		for k, expected := range tt.params {
			if got := ctx.Params[k]; got != expected {
				t.Errorf(
					"lookup(%q): expected param %s=%q, got %q",
					tt.path,
					k,
					expected,
					got,
				)
			}
		}
	}
}

func TestRadixNode_lookup_CatchAllAndParams(t *testing.T) {
	tree := &RadixNode{}

	tree.addRoute(
		"/orgs/:org/repos/:repo/branches/:branch/files/*filepath",
		300,
	)

	ctx, err := tree.lookup(
		"/orgs/openai/repos/chatgpt/branches/main/files/src/router/radix.go",
	)

	if err != nil {
		t.Fatal("Error occurred in lookup: " + err.Error())
	}

	if ctx.Params["org"] != "openai" {
		t.Errorf("Expected org=openai, got %q", ctx.Params["org"])
	}

	if ctx.Params["repo"] != "chatgpt" {
		t.Errorf("Expected repo=chatgpt, got %q", ctx.Params["repo"])
	}

	if ctx.Params["branch"] != "main" {
		t.Errorf("Expected branch=main, got %q", ctx.Params["branch"])
	}

	if ctx.Params["filepath"] != "/src/router/radix.go" {
		t.Errorf(
			"Expected filepath=/src/router/radix.go, got %q",
			ctx.Params["filepath"],
		)
	}

	if ctx.Data != 300 {
		t.Errorf("Expected data=300, got %v", ctx.Data)
	}
}

func TestRadixNode_lookup_MixedRoutes(t *testing.T) {
	tree := &RadixNode{}

	tree.addRoute("/", 1)
	tree.addRoute("/users/:id", 2)
	tree.addRoute("/users/:id/posts/:postID", 3)
	tree.addRoute("/static/*file", 4)
	tree.addRoute("/api/:version/docs/*page", 5)
	tree.addRoute("/orgs/:org/repos/:repo/branches/:branch/files/*filepath", 6)

	tests := []struct {
		path   string
		data   int
		params map[string]string
	}{
		{
			path:   "/",
			data:   1,
			params: map[string]string{},
		},
		{
			path: "/users/42",
			data: 2,
			params: map[string]string{
				"id": "42",
			},
		},
		{
			path: "/users/42/posts/99",
			data: 3,
			params: map[string]string{
				"id":     "42",
				"postID": "99",
			},
		},
		{
			path: "/static/css/main.css",
			data: 4,
			params: map[string]string{
				"file": "/css/main.css",
			},
		},
		{
			path: "/api/v1/docs/getting-started/install",
			data: 5,
			params: map[string]string{
				"version": "v1",
				"page":    "/getting-started/install",
			},
		},
		{
			path: "/orgs/openai/repos/chatgpt/branches/main/files/src/router/radix.go",
			data: 6,
			params: map[string]string{
				"org":      "openai",
				"repo":     "chatgpt",
				"branch":   "main",
				"filepath": "/src/router/radix.go",
			},
		},
	}

	for _, tt := range tests {
		ctx, err := tree.lookup(tt.path)
		if err != nil {
			t.Fatalf("lookup(%q) failed: %v", tt.path, err)
		}

		if ctx.Data != tt.data {
			t.Errorf(
				"lookup(%q): expected data=%d, got %v",
				tt.path,
				tt.data,
				ctx.Data,
			)
		}

		for k, expected := range tt.params {
			if got := ctx.Params[k]; got != expected {
				t.Errorf(
					"lookup(%q): expected %s=%q, got %q",
					tt.path,
					k,
					expected,
					got,
				)
			}
		}
	}
}
