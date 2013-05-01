// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package heroku

import (
	"errors"
	"github.com/jmcvetta/restclient"
	"net/url"
	"log"
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

type mapApp struct {
	App map[string]string `json:"app"`
}

func (h *Heroku) userinfo() *url.Userinfo {
	return url.UserPassword("", h.ApiKey)
}

// Apps queries Heroku for all applications owned by account, and returns a
// map keyed with application names.
func (h *Heroku) Apps() (map[string]*App, error) {
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
	m := make(map[string]*App, len(res))
	for _, a := range res {
		m[a.Name] = a
	}
	return m, nil
}

func (h *Heroku) NewApp(name, stack string) (*App, error) {
	url := h.ApiHref + "/apps"
	a := new(App)
	m := make(map[string]string)
	if name != "" {
		m["name"] = name
	}
	if stack != "" {
		m["stack"] = stack
	}
	payload := &mapApp{m}
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
		log.Println(status)
		log.Println(*e)
		return nil, BadResponse
	}
	return a, nil
}


// DestroyApp deletes an application from Heroku.
func (h *Heroku) DestroyApp(name string) error {
	url := h.ApiHref + "/apps/" + name
	e := new(interface{})
	rr := restclient.RequestResponse{
		Url: url,
		Method: "DELETE",
		Userinfo: h.userinfo(),
		Error: e,
	}
	status, err := h.rc.Do(&rr)
	if err != nil {
		return err
	}
	if status != 200 {
		log.Println(status)
		log.Println(*e)
		return BadResponse
	}
	return nil // Successful delete
}
