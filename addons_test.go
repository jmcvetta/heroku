// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package heroku

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestAddons(t *testing.T) {
	h := setup(t)
	defer cleanup(t, h)
	addons, err := h.Addons()
	if err != nil {
		t.Fatal(err)
	}
	prettyPrint(addons)
	assert.NotEqual(t, 0, len(addons))
}
