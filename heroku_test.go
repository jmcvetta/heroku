// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package heroku

import (
	"encoding/json"
	"github.com/bmizerany/assert"
	// "github.com/kr/pretty"
	"github.com/darkhelmet/env"
	"github.com/jmcvetta/randutil"
	"net/http"
	// "net/http/httptest"
	// "reflect"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"
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

var testApps = []*App{
	&App{
		Id:                1,
		Name:              "foo",
		CreateStatus:      time.Now().String(),
		CreatedAt:         time.Now().String(),
		Stack:             "cedar",
		RequestedStack:    "",
		RepoMigrateStatus: "complete",
		SlugSize:          2412544,
		RepoSize:          1777664,
		Dynos:             3,
		Workers:           1,
	},
	&App{
		Id:                2,
		Name:              "bar",
		CreateStatus:      time.Now().String(),
		CreatedAt:         time.Now().String(),
		Stack:             "cedar",
		RequestedStack:    "",
		RepoMigrateStatus: "complete",
		SlugSize:          1234,
		RepoSize:          5678,
		Dynos:             1,
		Workers:           0,
	},
}

func HandleGetApps(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	enc.Encode(&testApps)
}

func HandleNewApp(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()
	a0 := App{}
	dec.Decode(&a0)
	if a0.Name != "foo" {
		msg := fmt.Sprintf("Expected name='foo', got '%s'.", a0.Name)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	if a0.Stack != "bar" {
		msg := fmt.Sprintf("Expected stack='bar', got '%s'.", a0.Stack)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	enc := json.NewEncoder(w)
	a1 := testApps[0]
	enc.Encode(&a1)
	w.WriteHeader(http.StatusAccepted)
}

/*
func TestGetApps(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(HandleGetApps))
	defer srv.Close()
	h := setup(t, srv)
	m, err := h.Apps()
	if err != nil {
		t.Fatal(err)
	}
	for _, a0 := range testApps {
		a1, ok := m[a0.Id]
		if !ok {
			t.Fatal("Apps() failed to return app ", a0.Id)
		}
		assert.T(t, reflect.DeepEqual(a0, a1))
	}
}
*/

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
