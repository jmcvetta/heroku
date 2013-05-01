// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package heroku

import (
	"encoding/json"
	"github.com/bmizerany/assert"
	"github.com/darkhelmet/env"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func setup(t *testing.T) *Heroku {
	key := env.StringDefault("HEROKU_API_KEY", "")
	if key == "" {
		t.Fatal("HEROKU_API_KEY environment variable not set")
	}
	h := new(Heroku)
	h.ApiKey = key
	return h
}

var testApps = []App{
	App{
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
	App{
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

func TestGetApps(t *testing.T) {
	h := setup(t)
	srv := httptest.NewServer(http.HandlerFunc(HandleGetApps))
	defer srv.Close()
	m, err := h.Apps()
	if err != nil {
		t.Fatal(err)
	}
	for _, a0 := range m {
		a1, ok := m[a0.Id]
		if !ok {
			t.Error("Apps() failed to return app ", a0.Id)
		}
		assert.T(t, reflect.DeepEqual(a0, a1))
	}
}
