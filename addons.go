// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package heroku

import (
	"github.com/jmcvetta/restclient"
)

/*
HTTP/1.1 200 OK
[
  {
    "name": "example:basic",
    "description": "Example Basic",
    "url": "http://devcenter.heroku.com/articles/example-basic",
    "state": "public",
    "beta": false,
  }
]
*/

type Addon struct {
	Name        string
	Description string
	Url         string
	State       string
	Beta        string
}

func (h *Heroku) Addons() ([]*Addon, error) {
	url := HerokuApi + "/addons"
	res := []*Addon{}
	rr := restclient.RequestResponse{
		Url:      url,
		Userinfo: h.userinfo(),
		Result:   &res,
	}
	status, err := h.rc.Do(&rr)
	if status != 200 || err != nil {
		prettyPrint(err)
		return nil, err
	}
	return res, nil
}
