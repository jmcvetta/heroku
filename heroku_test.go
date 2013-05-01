// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package heroku

import (
	"github.com/darkhelmet/env"
	"github.com/kr/pretty"
	"log"
	"testing"
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

func TestGetApps(t *testing.T) {
	a := setup(t)
	apps, err := a.Apps()
	if err != nil {
		t.Fatal(err)
	}
	for _, app := range apps {
		log.Println(pretty.Sprintf("%# v", app))
	}
}
