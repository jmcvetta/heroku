// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package heroku

import (
	"errors"
	"github.com/jmcvetta/restclient"
	"net/url"
)

const HerokuApi = "https://api.heroku.com"

var (
	BadResponse = errors.New("Bad response from server")
)

func NewHeroku(apiKey string) *Heroku {
	h := Heroku{
		ApiKey:  apiKey,
		ApiHref: HerokuApi,
		rc:      restclient.New(),
	}
	return &h
}

type Heroku struct {
	ApiKey  string
	ApiHref string
	rc      *restclient.Client
}

func (h *Heroku) userinfo() *url.Userinfo {
	return url.UserPassword("", h.ApiKey)
}

