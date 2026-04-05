package assets

import (
	"strings"
)

var allowedDomains = []string{
	"www.bing.com",
	"th.bing.com",
	"image-assets.soutushenqi.com",
	"i0.hdslb.com",
	"c-ssl.duitang.com",
	"gd-hbimg.huaban.com",
	"pic.rmb.bdstatic.com",
}

func IsAllowedHost(host string) bool {
	for _, allowed := range allowedDomains {
		if host == allowed || strings.HasSuffix(host, "."+allowed) {
			return true
		}
	}
	return false
}
