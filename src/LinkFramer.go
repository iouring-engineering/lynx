package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

func anyToString(value any) string {
	switch v := value.(type) {
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func frameCompleteUrl(linkData DbShortLink) string {
	var m map[string]any = make(map[string]any)
	err := json.Unmarshal([]byte(linkData.Data), &m)
	ErrorLogger.Println(err)
	m["shortcode"] = linkData.ShortCode
	var urlData = ""
	var idx = 0
	for key, value := range m {
		urlData += fmt.Sprintf("%s=%s", url.QueryEscape(key), url.QueryEscape(anyToString(value)))
		if idx != (len(m) - 1) {
			urlData += "&"
		}
		idx++
	}
	parsed, err := url.Parse(linkData.WebUrl)
	if len(parsed.Query()) > 0 {
		if linkData.WebUrl == "" {
			return fmt.Sprintf("%s&%s", config.AppConfig.DefaultFallbackUrl, urlData)
		}
		return fmt.Sprintf("%s&%s", linkData.WebUrl, urlData)
	}

	if linkData.WebUrl == "" {
		return fmt.Sprintf("%s?%s", config.AppConfig.DefaultFallbackUrl, urlData)
	}
	return fmt.Sprintf("%s?%s", linkData.WebUrl, urlData)
}

func frameAndroidUrl(android, shortCode string) string {
	var parsed MobileInputs
	json.Unmarshal([]byte(android), &parsed)
	if parsed.Fbl == "" || android == "" {
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
	if parsed.Fbl == "" || ios == "" {
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
