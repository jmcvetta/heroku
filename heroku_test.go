// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package heroku

import (
	"encoding/json"
	"github.com/darkhelmet/env"
	"github.com/kr/pretty"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func setup(t *testing.T) *Account {
	key := env.StringDefault("HEROKU_API_KEY", "")
	if key == "" {
		t.Fatal("HEROKU_API_KEY environment variable not set")
	}
	a := new(Account)
	a.ApiKey = key
	return a
}

func HandleGetApps(w http.ResponseWriter, r *http.Request) {
	a0 := App{
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
	}
	a1 := App{
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
	}
	enc := json.NewEncoder(w)
	resp := []App{a0, a1}
	enc.Encode(&resp)
}

func TestGetApps(t *testing.T) {
	a := setup(t)
	srv := httptest.NewServer(http.HandlerFunc(HandleGetApps))
	defer srv.Close()
	apps, err := a.Apps()
	if err != nil {
		t.Fatal(err)
	}
	for _, app := range apps {
		log.Println(pretty.Sprintf("%# v", app))
	}
}
