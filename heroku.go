// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package heroku

import (
	"errors"
	"github.com/jmcvetta/restclient"
	"net/url"
)

const HerokuApi   = "https://api.heroku.com"

var (
	BadResponse = errors.New("Bad response from server")
)

func NewHeroku(apiKey string) *Heroku {
	h := Heroku{
		ApiKey: apiKey,
		ApiHref: HerokuApi,
		rc: restclient.New(),
	}
	return &h
}

type Heroku struct {
	ApiKey string
	ApiHref string
	rc *restclient.Client
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

type mss map[string]string


func (h *Heroku) userinfo() *url.Userinfo {
	return url.UserPassword("", h.ApiKey)
}

// Apps queries Heroku for all applications owned by account, and returns a
// map keyed with app IDs.
func (h *Heroku) Apps() (map[int64]*App, error) {
	url := h.ApiHref + "/apps"
	res := []*App{}
	e := new(interface{})
	rr := restclient.RequestResponse{
		Url:      url,
		Method:   "GET",
		Userinfo: h.userinfo(),
		Result:   &res,
		Error: e,
	}
	status, err := h.rc.Do(&rr)
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

func (h *Heroku) NewApp(name, stack string) (*App, error) {
	url := h.ApiHref + "/apps"
	a := new(App)
	payload := struct {
		App mss `json:"app"`
	}{
		App: mss{"name": name, "stack": stack,},
	}
	e := new(interface{})
	rr := restclient.RequestResponse{
		Url: url,
		Method: "POST",
		Userinfo: h.userinfo(),
		Result: a,
		Data: payload,
		Error: e,
	}
	status, err := h.rc.Do(&rr)
	if err != nil {
		return nil, err
	}
	if status != 202 {
		return nil, BadResponse
	}
	return a, nil
}
