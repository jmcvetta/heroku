// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package heroku

import (
	"encoding/json"
	"fmt"
	"github.com/jmcvetta/restclient"
	"log"
)

type Addon struct {
	Name        string
	Description *string
	Url         *string
	Beta        *bool
	State       *string
	Attachable  *json.RawMessage
	Price       *struct {
		Cents int
		Unit  string
	}
	Slug             *string
	Tos              *bool   `json:"terms_of_service"`
	DynoHours        *bool   `json:"consumes_dyno_hours"`
	PlanDescription  *string `json:"plan_description"`
	GroupDescription *string `json:"group_description"`
	Selective        *json.RawMessage
	Configured       *bool
}

type AddonStatus struct {
	Status  string
	Message string
	Price   string
}

// Addons lists all addons that can be installed.
func (h *Heroku) Addons() ([]*Addon, error) {
	url := HerokuApi + "/addons"
	res := []*Addon{}
	rr := restclient.RequestResponse{
		Url:      url,
		Method:   "GET",
		Userinfo: h.userinfo(),
		Result:   &res,
	}
	status, err := h.rc.Do(&rr)
	if status != 200 || err != nil {
		log.Println(err)
		return nil, err
	}
	return res, nil
}

// InstallAddon provisions an app with an addon.
func (h *Heroku) InstallAddon(app, addon string) (*AddonStatus, error) {
	url := HerokuApi + fmt.Sprintf("/apps/%s/addons/%s", app, addon)
	res := AddonStatus{}
	e := new(HerokuError)
	rr := restclient.RequestResponse{
		Url:      url,
		Method:   "POST",
		Userinfo: h.userinfo(),
		Error:    e,
		Result:   &res,
	}
	status, err := h.rc.Do(&rr)
	if err != nil {
		log.Println(err)
		prettyPrint(e)
		return nil, err
	}
	if status != 200 {
		log.Println(status)
		prettyPrint(rr.Error)
		return nil, BadResponse
	}
	return &res, nil

}

// AppAddons lists the addons with which an app is provisioned.
func (h *Heroku) AppAddons(app string) ([]*Addon, error) {
	url := HerokuApi + fmt.Sprintf("/apps/%s/addons", app)
	res := []*Addon{}
	rr := restclient.RequestResponse{
		Url:      url,
		Method:   "GET",
		Userinfo: h.userinfo(),
		Result:   &res,
	}
	status, err := h.rc.Do(&rr)
	if status != 200 || err != nil {
		log.Println(err)
		return nil, err
	}
	return res, nil
}

/*

This part of the Heroku API doesn't seem to work right.  Even a query
formed with the example generator on the API docs page does not perform
as expected.

// UpgradeAddon changes the plan type for an installed addon.
func (h *Heroku) UpgradeAddon(app, addon string) (*AddonStatus, error) {
	url := HerokuApi + fmt.Sprintf("/apps/%s/addons/%s", app, addon)
	res := AddonStatus{}
	e := new(HerokuError)
	rr := restclient.RequestResponse{
		Url:      url,
		Method:   "PUT",
		Userinfo: h.userinfo(),
		Result:   &res,
		Error:    e,
		Data: nil,
	}
	status, err := h.rc.Do(&rr)
	if err != nil {
		prettyPrint(err)
		prettyPrint(e)
		return nil, err
	}
	if status != 200 {
		log.Println(status)
		prettyPrint(rr.Error)
		return nil, BadResponse
	}
	return &res, nil
}

*/
