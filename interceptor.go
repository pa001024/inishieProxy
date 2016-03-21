package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"path"
	"regexp"

	"github.com/pa001024/reflex/util"
)

// 拦截器
type ResourceInterceptor struct {
	config    *RemapConfig
	actionMap map[string][]TranlatePoint // 页面路径->注入点
	pointMap  map[string]Filter          // 注入点名字->注入器
}

func NewResourceInterceptor(dir, lang string) (rst *ResourceInterceptor) {
	ad := path.Join(dir, lang)
	rst = &ResourceInterceptor{
		actionMap: make(map[string][]TranlatePoint),
		pointMap:  make(map[string]Filter),
	}
	_, err := os.Stat(ad)
	util.Try(err)

	rst.loadConfig(dir)
	rst.loadPoints(ad)
	return
}

type RemapConfig struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	GameVersion string          `json:"game-version"`
	Last        string          `json:"last"`
	Points      []TranlatePoint `json:"points"`
}
type TranlatePoint struct {
	File   string `json:"file"`
	Name   string `json:"name"`
	Action string `json:"action"`
	Key    string `json:"key"`
	Mode   string `json:"mode"`
}

// 从配置文件加载拦截器配置
func (this *ResourceInterceptor) loadConfig(dir string) *RemapConfig {
	// load config
	fin, err := os.Open(path.Join(dir, "remap.json"))
	util.Try(err)
	d := json.NewDecoder(fin)
	remap := &RemapConfig{}
	d.Decode(remap)
	fin.Close()
	this.config = remap
	return remap
}

// 初始化拦截器
func (this *ResourceInterceptor) loadPoints(ad string) {
	for _, v := range this.config.Points {
		regexp.MustCompile(v.Key) // check the regexp is illegal
		fn := path.Join(ad, v.File)
		_, err := os.Stat(fn)
		if err != nil && os.IsNotExist(err) {
			util.WARN.Log("[Missing Translate File] ", fn)
			continue
		}
		this.pointMap[v.Name] = NewFilter(v.Mode, fn)
		if _, ok := this.actionMap[v.Action]; !ok {
			this.actionMap[v.Action] = make([]TranlatePoint, 0, 3)
		}
		this.actionMap[v.Action] = append(this.actionMap[v.Action], v)
		util.DEBUG.Log("[Interceptor] Load ", v.Action, " -> ", v.Name)
	}
}

func (this *ResourceInterceptor) Process(req *http.Request, query string, buf *bytes.Buffer) {
	util.DEBUG.Log("[Interceptor] matching ", req.URL.Path)
	if vls, ok := this.actionMap[req.URL.Path]; ok {
		for _, v := range vls {
			if p, ok := this.pointMap[v.Name]; ok {
				if p == nil {
					util.DEBUG.Log("[Interceptor] ", v.Name, " p == nil")
					continue
				}
				util.DEBUG.Log("[Interceptor] ", v.Name, v.Key)
				if (v.Key == "-" && query == "") || v.Key == "*" || v.Key == query || regexp.MustCompile(v.Key).MatchString(query) {
					s := this.pointMap[v.Name].ReplaceAll(buf.String())
					buf.Reset()
					buf.WriteString(s)
				}
			}
		}
	}
}
