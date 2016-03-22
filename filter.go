package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/pa001024/reflex/util"
)

type Filter interface {
	ReplaceAll(s string) string
}

func NewFilter(mode, filename string) Filter {
	switch mode {
	case "ex-text-replace":
		return NewRexKVRplFilter(filename)
	case "text-replace":
		return NewTextRplFilter(filename)
	case "file-replace":
		return NewFileRplFilter(filename)
	case "regexp-replace":
		return NewRegexpRplFilter(filename)
	}
	return nil
}

// 强模式KV文本过滤器
type RexKVRplFilter struct {
	Filter
	dict  map[string]string
	ex    *regexp.Regexp
	preex *regexp.Regexp
	subex *regexp.Regexp
}

type RexKVConfig struct {
	RegexpStr string `json:"regexp_str"`
	PrefixStr string `json:"prefix_str"`
	SubfixStr string `json:"subfix_str"`
	KvFile    string `json:"kv_file"`
}

func NewRexKVRplFilter(filename string) (v *RexKVRplFilter) {
	fin, err := os.Open(filename)
	ad := path.Dir(filename)
	util.Try(err)
	config := &RexKVConfig{}
	de := json.NewDecoder(fin)
	de.Decode(config)
	fin.Close()

	bin, err := ioutil.ReadFile(path.Join(ad, config.KvFile))
	util.Try(err)
	ex := regexp.MustCompile(`\r?\n(?:\s*\r?\n)*|=`)
	source := ex.Split(string(bin), -1)

	v = &RexKVRplFilter{dict: make(map[string]string),
		ex:    regexp.MustCompile(config.RegexpStr),
		preex: regexp.MustCompile("^" + config.PrefixStr),
		subex: regexp.MustCompile(config.SubfixStr + "$"),
	}

	for i := 0; i < len(source); i += 2 {
		v.dict[source[i]] = source[i+1]
	}
	return
}

func (this *RexKVRplFilter) ReplaceAll(s string) string {
	return this.ex.ReplaceAllStringFunc(s, func(m string) string {
		pre := this.preex.FindString(m)
		sub := this.subex.FindString(m)
		mid := m[len(pre) : len(m)-len(sub)]
		if v, ok := this.dict[mid]; ok {
			return pre + v + sub
		}
		return m
	})
}

// 多文本替换过滤器
type TextRplFilter struct {
	Filter
	Source   []string
	replacer *strings.Replacer
}

func NewTextRplFilter(filename string) *TextRplFilter {
	bin, err := ioutil.ReadFile(filename)
	util.Try(err)
	ex := regexp.MustCompile(`\r?\n(?:\s*\r?\n)*|=`)
	source := ex.Split(string(bin), -1)
	if len(source)%2 == 1 {
		source = source[:len(source)-1]
	}
	for i, v := range source {
		util.DEBUG.Log(i, v)
	}
	return &TextRplFilter{
		Source:   source,
		replacer: strings.NewReplacer(source...),
	}
}
func (this *TextRplFilter) ReplaceAll(s string) string {
	return this.replacer.Replace(s)
}

// 文件替换器
type FileRplFilter struct {
	Filter
	filename string
}

func (this *FileRplFilter) ReplaceAll(s string) string {
	bin, _ := ioutil.ReadFile(this.filename)
	return string(bin)
}

func NewFileRplFilter(filename string) *FileRplFilter {
	_, err := os.Stat(filename)
	util.Try(err)
	return &FileRplFilter{
		filename: filename,
	}
}

// 正则表达式过滤器
type RegexpRplFilter struct {
	Filter
	source []string
}

func (this *RegexpRplFilter) ReplaceAll(s string) string {
	for i := 0; i < len(this.source); i += 2 {
		ex := regexp.MustCompile(this.source[i])
		s = ex.ReplaceAllString(s, this.source[i+1])
	}
	return s
}

func NewRegexpRplFilter(filename string) *RegexpRplFilter {
	fin, err := os.Open(filename)
	util.Try(err)
	defer fin.Close()
	de := json.NewDecoder(fin)
	var source []string
	de.Decode(&source)
	return &RegexpRplFilter{
		source: source,
	}
}
