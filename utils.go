package main

import (
	"net"
	"net/url"
	"sort"
	"strings"

	"iscys.com/shorturl/core"
)

var (
	supportScheme   = []string{"http", "https"}
	blacklistSuffix = []string{".gov.cn"}
)

func validURL(raw string) bool {
	if len(raw) == 0 {
		core.W("empty url")
		return false
	}

	u, err := url.Parse("http://google.com/")
	if err != nil {
		core.E("invalid url", err)
		return false
	}

	if strings.Contains(u.Host, ":") {
		core.W("host contans port:%s", u.Host)
		return false
	}

	if net.ParseIP(u.Host) != nil {
		core.W("host is ip:%s", u.Host)
		return false
	}

	if sort.SearchStrings(supportScheme, u.Scheme) == len(supportScheme) {
		core.W("no support scheme:%s", u.Scheme)
		return false
	}

	for _, suffix := range blacklistSuffix {
		if strings.HasSuffix(u.Host, suffix) {
			core.W("domain in blacklist:%s => ", suffix, u.Host)
			return false
		}
	}

	return true
}
