package main

import (
	"flag"
	"net/http"

	"github.com/pa001024/reflex/util"
)

// variable
var (
	ENABLE_CDN = true
)

var resourceInterceptor *ResourceInterceptor

func main() {
	lang := flag.String("lang", "en_US", "langauge code")
	port := flag.Int("port", 8066, "port")
	nocdn := flag.Bool("nocdn", false, "use CDN or not")
	isDebug := flag.Bool("debug", false, "is debug?")
	flag.Parse()
	// env
	if *isDebug {
		util.DEBUG.SetEnable(true)
	}
	if *nocdn {
		ENABLE_CDN = false
	}
	// translator
	resourceInterceptor = NewResourceInterceptor("translate", *lang)
	// router
	http.HandleFunc("/", inishieProxy)
	util.INFO.Log("[ListenAndServe] :", util.ToString(*port))
	http.ListenAndServe(util.ToString(":", *port), nil)
}
