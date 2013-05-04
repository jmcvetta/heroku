// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package heroku

import (
	"github.com/bmizerany/assert"
	"testing"
	"time"
)

const pgFree = "heroku-postgresql:dev"

func TestAddons(t *testing.T) {
	h := setup(t)
	defer cleanup(t, h)
	addons, err := h.Addons()
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEqual(t, 0, len(addons))
}

func TestInstallAddon(t *testing.T) {
	h := setup(t)
	defer cleanup(t, h)
	a, _ := h.NewApp("", "")
	aStat, err := a.InstallAddon(pgFree)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEqual(t, nil, aStat.Message)
	assert.Equal(t, "free", aStat.Price)
}

func TestAppAddons(t *testing.T) {
	h := setup(t)
	defer cleanup(t, h)
	a, _ := h.NewApp("", "")
	a.InstallAddon(pgFree)
	// Sleep while provisioning occurs.
	dur, _ := time.ParseDuration("2s")
	time.Sleep(dur)
	addons, err := a.Addons()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(addons))
	assert.Equal(t, pgFree, addons[0].Name)
}
