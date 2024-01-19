package main

import (
	"encoding/json"
	"fmt"
	"net/url"
)

func frameCompleteUrl(linkData DbShortLink) string {
	unescaped := url.QueryEscape(linkData.Data)
	if linkData.WebUrl == "" {
		return fmt.Sprintf("%s?data=%s", config.AppConfig.DefaultFallbackUrl, unescaped)
	}
	return fmt.Sprintf("%s?data=%s", linkData.WebUrl, unescaped)
}

func frameAndroidUrl(android, shortCode string) string {
	var parsed MobileInputs
	json.Unmarshal([]byte(android), &parsed)
	if parsed.Fbl == "" {
		if config.AppConfig.Android.Behaviour == APP_SEARCH {
			var play = config.AppConfig.Android.GooglePlaySearchUrl
			return fmt.Sprintf("%s&referrer=%s", play, shortCode)
		} else if config.AppConfig.Android.AndroidDefaultWebUrl != "" {
			return config.AppConfig.Android.AndroidDefaultWebUrl
		} else {
			return config.AppConfig.DefaultFallbackUrl
		}
	}
	return parsed.Fbl
}

func frameIosUrl(ios string) string {
	var parsed MobileInputs
	json.Unmarshal([]byte(ios), &parsed)
	if parsed.Fbl == "" {
		if config.AppConfig.Ios.Behaviour == APP_SEARCH {
			return config.AppConfig.Ios.AppStoreSearchUrl
		} else if config.AppConfig.Ios.IosDefaultWebUrl != "" {
			return config.AppConfig.Ios.IosDefaultWebUrl
		} else {
			return config.AppConfig.DefaultFallbackUrl
		}
	}
	return parsed.Fbl
}
