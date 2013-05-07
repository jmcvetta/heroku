// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package heroku

import (
	"github.com/darkhelmet/env"
	"github.com/jmcvetta/randutil"
	"log"
	"net/http"
	"strings"
	"testing"
)

func setup(t *testing.T) *Heroku {
	log.SetFlags(log.Lshortfile)
	key := env.String("HEROKU_TEST_API_KEY")
	h := NewHeroku(key)
	h.rc.Log = false
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

func TestPrettyPrint(t *testing.T) {
	// This test can only fail with a panic
	prettyPrint(map[string]string{"foo": "bar"})
}
