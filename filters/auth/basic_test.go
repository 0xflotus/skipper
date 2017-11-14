package auth

import (
	"github.com/zalando/skipper/filters/filtertest"
	"net/http"
	"testing"
)

type createTestItem struct {
	msg  string
	args []interface{}
	err  bool
}

func TestWithMissingAuth(t *testing.T) {
	spec := NewBasicAuth()
	f, err := spec.CreateFilter([]interface{}{"testdata/htpasswd"})
	if err != nil {
		t.Error(err)
		return
	}

	req, err := http.NewRequest("GET", "https://www.example.org/", nil)
	if err != nil {
		t.Error(err)
		return
	}

	expectedBasicAuthHeaderValue := ForceBasicAuthHeaderValue + `"Basic Realm"`

	ctx := &filtertest.Context{FRequest: req}
	f.Request(ctx)
	if ctx.Response().Header.Get(ForceBasicAuthHeaderName) != expectedBasicAuthHeaderValue && ctx.Response().StatusCode == 401 && ctx.Served() {
		t.Error("Authentication header wrong/missing")
	}
}

func TestWithWrongAuth(t *testing.T) {
	spec := NewBasicAuth()
	f, err := spec.CreateFilter([]interface{}{"testdata/htpasswd", "My Website"})
	if err != nil {
		t.Error(err)
		return
	}

	req, err := http.NewRequest("GET", "https://www.example.org/", nil)
	req.SetBasicAuth("myName", "wrongPassword")
	if err != nil {
		t.Error(err)
		return
	}

	expectedBasicAuthHeaderValue := ForceBasicAuthHeaderValue + `"My Website"`

	ctx := &filtertest.Context{FRequest: req}
	f.Request(ctx)
	if ctx.Response().Header.Get(ForceBasicAuthHeaderName) != expectedBasicAuthHeaderValue && ctx.Response().StatusCode == 401 && ctx.Served() {
		t.Error("Authentication header wrong/missing")
	}
}

func TestWithSuccessfulAuth(t *testing.T) {
	spec := NewBasicAuth()
	f, err := spec.CreateFilter([]interface{}{"testdata/htpasswd"})
	if err != nil {
		t.Error(err)
		return
	}

	req, err := http.NewRequest("GET", "https://www.example.org/", nil)
	req.SetBasicAuth("myName", "myPassword")
	if err != nil {
		t.Error(err)
		return
	}

	ctx := &filtertest.Context{FRequest: req}
	f.Request(ctx)
	if ctx.Served() && ctx.Response().StatusCode != 401 {
		t.Error("Authentication not successful")
	}
}
