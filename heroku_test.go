// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package heroku

import (
	"github.com/bmizerany/assert"
	"github.com/darkhelmet/env"
	"github.com/jmcvetta/randutil"
	"log"
	"net/http"
	"strings"
	"testing"
)

func setup(t *testing.T) *Heroku {
	log.SetFlags(log.Lshortfile)
	key := env.StringDefault("HEROKU_API_KEY", "")
	if key == "" {
		t.Fatal("HEROKU_API_KEY environment variable not set")
	}
	h := NewHeroku(key)
	// addr := "http://" + s.Listener.Addr().String()
	// h.ApiHref = addr
	// h.rc.UnsafeBasicAuth = true
	h.rc.LogRequestResponse = false
	return h
}

func BadResponder(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "", 999)
}

func appName(t *testing.T) string {
	rnd, err := randutil.AlphaString(25)
	if err != nil {
		t.Fatal(err)
	}
	rnd = strings.ToLower(rnd)
	return "test-" + rnd
}

func cleanup(t *testing.T, h *Heroku) {
	m, err := h.Apps()
	if err != nil {
		t.Fatal(err)
	}
	for name := range m {
		err = h.DestroyApp(name)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestNewApp(t *testing.T) {
	h := setup(t)
	defer cleanup(t, h)
	a0name := appName(t)
	a0, err := h.NewApp(a0name, "")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, a0name, a0.Name)
	//
	// Duplicate name should cause error
	//
	_, err = h.NewApp(a0name, "")
	assert.NotEqual(t, nil, err)
}

func TestGetApp(t *testing.T) {
	h := setup(t)
	defer cleanup(t, h)
	a0name := appName(t)
	h.NewApp(a0name, "")
	a0, err := h.App(a0name)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, a0name, a0.Name)
}

func TestDestroyApp(t *testing.T) {
	h := setup(t)
	defer cleanup(t, h)
	a0name := appName(t)
	a0, _ := h.NewApp(a0name, "")
	err := h.DestroyApp(a0.Name)
	if err != nil {
		t.Fatal(err)
	}
	m, _ := h.Apps()
	_, ok := m[a0name]
	if ok {
		t.Error("Failed to delete app ", a0name)
	}
}

func TestGetApps(t *testing.T) {
	h := setup(t)
	defer cleanup(t, h)
	a0name := appName(t)
	a1name := appName(t)
	h.NewApp(a0name, "")
	h.NewApp(a1name, "")
	m, err := h.Apps()
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{a0name, a1name} {
		_, ok := m[name]
		if !ok {
			t.Error("Could not retrieve app ", name)
		}
	}
}

func TestMaintenanceMode(t *testing.T) {
	h := setup(t)
	defer cleanup(t, h)
	a0name := appName(t)
	a0, _ := h.NewApp(a0name, "")
	err := h.MaintenanceMode(a0.Name, true)
	if err != nil {
		t.Fatal(err)
	}
	err = h.MaintenanceMode(a0.Name, false)
	if err != nil {
		t.Fatal(err)
	}
}
