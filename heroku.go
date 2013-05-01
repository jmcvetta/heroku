// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package heroku

import (
	"errors"
	"github.com/jmcvetta/restclient"
	"net/url"
)

var (
	HerokuApi   = "https://api.heroku.com"
	BadResponse = errors.New("Bad response from server")
)

type Heroku struct {
	ApiKey string
}

// An App is a Heroku application.
type App struct {
	Id                int64
	Name              string
	CreateStatus      string `json:"create_status"`
	CreatedAt         string `json:"created_at"`
	Stack             string
	RequestedStack    string `json:"requested_stack"`
	RepoMigrateStatus string `json:"repo_migrate_status"`
	SlugSize          int    `json:"slug_size"`
	RepoSize          int    `json:"repo_size"`
	Dynos             int
	Workers           int
}

func (h *Heroku) userinfo() *url.Userinfo {
	return url.UserPassword("", h.ApiKey)
}

// Apps queries Heroku for all applications owned by account, and returns a
// map keyed with app IDs.
func (h *Heroku) Apps() (map[int64]*App, error) {
	url := HerokuApi + "/apps"
	res := []*App{}
	rr := restclient.RequestResponse{
		Url:      url,
		Method:   "GET",
		Userinfo: h.userinfo(),
		Result:   &res,
	}
	status, err := restclient.Do(&rr)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, BadResponse
	}
	m := make(map[int64]*App, len(res))
	for _, a := range res {
		m[a.Id] = a
	}
	return m, nil
}
