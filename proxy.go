package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/pa001024/reflex/util"
)

const (
	UPSTREAM_HOST = "inishie-dungeon.com"
	UPSTREAM_IP   = "203.138.99.138"
	UPSTREAM      = "http://" + UPSTREAM_HOST
)

// fetch from original server
func proxifyRequest(r *http.Request) (nr *http.Request, query string) {
	nu := *r.URL
	nu.Host = UPSTREAM_HOST
	nu.Scheme = "http"
	if r.Method == "GET" {
		nr, _ = http.NewRequest(r.Method, nu.String(), nil)
	} else {
		bin, _ := ioutil.ReadAll(r.Body)
		query, _ = url.QueryUnescape(string(bin))
		nr, _ = http.NewRequest(r.Method, nu.String(), bytes.NewBuffer(bin))

	}
	nr.Header = r.Header
	return
}

// POST
func inishieProxy(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" && r.URL.Path != "/" && path.Ext(r.URL.Path) != ".cgi" {
		fileProxy(w, r)
		return
	}
	req, query := proxifyRequest(r)

	util.DEBUG.Log("[inishieProxy] ", r.RequestURI, util.Sw(query == "", "", ": "+query))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		util.WARN.Log("[Proxy failed] ", r.RequestURI, ":\n\t\t\t", err)
	} else if res != nil && res.Body != nil {
		for i, v := range res.Trailer {
			for _, vv := range v {
				w.Header().Add(i, vv)
			}
		}
		bin, _ := ioutil.ReadAll(res.Body)
		buf := bytes.NewBuffer(bin)
		resourceInterceptor.Process(r, query, buf)
		io.Copy(w, buf)
	}
}

// Static File
func fileProxy(w http.ResponseWriter, r *http.Request) {
	_, err := os.Stat("dist" + r.URL.Path)
	if os.IsNotExist(err) {
		util.WARN.Log("[Missing File] ", UPSTREAM+r.URL.Path)
		rf, _ := http.Get(UPSTREAM + r.URL.Path)
		if rf != nil {
			os.MkdirAll("dist"+path.Dir(r.URL.Path), 0766)
			of, _ := os.OpenFile("dist"+r.URL.Path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
			io.Copy(of, rf.Body)
			rf.Body.Close()
		}
	}
	// 302 CDN
	if ENABLE_CDN && r.URL.Path != "/crossdomain.xml" && TestCDN(r.URL.Path) {
		RedirectToCDN(w, r, r.URL.Path)
	} else {
		util.DEBUG.Log("[fileProxy] ", r.URL.Path)
		http.ServeFile(w, r, "dist"+r.URL.Path)
	}
}
