// Copyright (c) 2013 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

package heroku

import (
	"encoding/json"
	"errors"
	"github.com/jmcvetta/restclient"
	"log"
	"net/url"
	"path/filepath"
	"runtime"
	"strconv"
)

const HerokuApi = "https://api.heroku.com"

var (
	BadResponse = errors.New("Bad response from server")
)

type HerokuError struct {
	Id    string
	Error string
}

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

func prettyPrint(v interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	lineNo := strconv.Itoa(line)
	file = filepath.Base(file)
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		log.Panic(err)
	}
	s := file + ":" + lineNo + ": \n" + string(b) + "\n"
	println(s)
}
