package main

import (
	"fmt"
	"net/url"
)

func frameBrowserUrl(linkData DbShortLink) string {
	if linkData.WebUrl == "" {
		return fmt.Sprintf("%s?data=%s", config.AppConfig.DefaultUrl, url.QueryEscape(linkData.Data))
	}
	return fmt.Sprintf("%s?data=%s", linkData.WebUrl, url.QueryEscape(linkData.Data))
}
