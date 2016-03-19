package main

import (
	"net/http"

	"github.com/pa001024/reflex/util"
)

const (
	CDN_ADDR = "http://inidunres.0w0.be"
)

var cdn_stat_cache map[string]bool = make(map[string]bool)

// POST
func TestCDN(path string) bool {
	if cdn_stat_cache[path] {
		return true
	}
	res, err := http.Head(CDN_ADDR + path)
	res.Body.Close()
	if err == nil && res.StatusCode == 200 {
		util.DEBUG.Log("[CDN] ", path)
		cdn_stat_cache[path] = true
		return true
	}
	util.WARN.Log("[CDN Missing] ", path)
	return false
}

func RedirectToCDN(w http.ResponseWriter, r *http.Request, path string) {
	http.Redirect(w, r, CDN_ADDR+path, 302) // TODO: change to 301
}
