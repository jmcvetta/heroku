// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package heroku

import (
	"github.com/jmcvetta/restclient"
	"log"
)

// An App is a Heroku application.
type App struct {
	Id                int64
	Name              string
	Dynos             *int
	Workers           *int
	RepoSize          *int `json:"repo_size"`
	SlugSize          *int `json:"slug_size"`
	Stack             *string
	RequestedStack    *string `json:"requested_stack"`
	CreateStatus      *string `json:"create_status"`
	RepoMigrateStatus *string `json:"repo_migrate_status"`
	OwnerEmail        *string `json:"owner_email"`
	OwnerName         *string `json:"owner_name"`
	DomainName        *struct {
		Id         *int64
		AppId      *int64 `json:"app_id"`
		Domain     *string
		BaseDomain *string `json:"base_domain"`
		Default    *bool
		CreatedAt  *string `json:"created_at"`
		UpdatedAt  *string `json:"updated_at"`
	} `json:"domain_name"`
	WebUrl    *string `json:"web_url"`
	GitUrl    *string `json:"git_url"`
	Tier      *string
	Region    *string
	CreatedAt *string `json:"created_at"`
	UpdatedAt *string `json:"updated_at"`
	h         *Heroku
}

type mapApp struct {
	App map[string]string `json:"app"`
}

// Apps lists all applications, returning a map keyed with app names.
func (h *Heroku) Apps() (map[string]*App, error) {
	url := h.ApiHref + "/apps"
	res := []*App{}
	e := new(interface{})
	rr := restclient.RequestResponse{
		Url:      url,
		Method:   "GET",
		Userinfo: h.userinfo(),
		Result:   &res,
		Error:    e,
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

// NewApp creates a new application.
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
		Url:      url,
		Method:   "POST",
		Userinfo: h.userinfo(),
		Result:   a,
		Data:     payload,
		Error:    e,
	}
	status, err := h.rc.Do(&rr)
	if err != nil {
		return nil, err
	}
	if status != 202 {
		log.Println(status)
		log.Println(*e)
		log.Println("name: ", name)
		log.Println("stack: ", stack)
		return nil, BadResponse
	}
	a.h = h
	return a, nil
}

// DestroyApp deletes an application.
func (h *Heroku) DestroyApp(appName string) error {
	url := h.ApiHref + "/apps/" + appName
	e := new(interface{})
	rr := restclient.RequestResponse{
		Url:      url,
		Method:   "DELETE",
		Userinfo: h.userinfo(),
		Error:    e,
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
	return nil
}

// App gets information about an application.
func (h *Heroku) App(appName string) (*App, error) {
	url := h.ApiHref + "/apps/" + appName
	e := new(interface{})
	a := new(App)
	rr := restclient.RequestResponse{
		Url:      url,
		Method:   "GET",
		Userinfo: h.userinfo(),
		Error:    e,
		Result:   a,
	}
	status, err := h.rc.Do(&rr)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		log.Println(status)
		log.Println(*e)
		return nil, BadResponse
	}
	a.h = h
	return a, nil

}

// MaintenanceMode toggles maintenance mode on an application.
func (h *Heroku) MaintenanceMode(appName string, modeOn bool) error {
	url := h.ApiHref + "/apps/" + appName + "/server/maintenance"
	e := new(interface{})
	payload := map[string]interface{}{}
	payload["app"] = appName
	payload["maintenance_mode"] = modeOn
	rr := restclient.RequestResponse{
		Url:      url,
		Method:   "POST",
		Userinfo: h.userinfo(),
		Error:    e,
		Data:     payload,
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
	return nil
}

// InstallAddon provisions this application with an addon.
func (a *App) InstallAddon(addon string) (*AddonStatus, error) {
	return a.h.InstallAddon(a.Name, addon)
}

// Addons lists the addons with which this app is provisioned.
func (a *App) Addons() ([]*Addon, error) {
	return a.h.AppAddons(a.Name)
}

/*
// UpgradeAddon changes the plan type of an installed addon.
func (a *App) UpgradeAddon(addon string) (*AddonStatus, error) {
	return a.h.UpgradeAddon(a.Name, addon)
}
*/
