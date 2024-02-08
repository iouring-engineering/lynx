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

func frameCompleteUrl(linkData DbShortLink, utm map[string]string) string {
	var m map[string]any = make(map[string]any)
	err := json.Unmarshal([]byte(linkData.Data), &m)
	if err != nil {
		ErrorLogger.Println(err)
	}
	var utmData = url.Values{}
	for key, value := range utm {
		utmData.Add(url.QueryEscape(key), url.QueryEscape(anyToString(value)))
	}
	if len(utm) > 0 {
		m["utm"] = utmData.Encode()
	}
	m["shortcode"] = linkData.ShortCode
	var urlData = url.Values{}
	for key, value := range m {
		urlData.Add(url.QueryEscape(key), url.QueryEscape(anyToString(value)))
	}
	parsed, err := url.Parse(linkData.WebUrl)
	if len(parsed.Query()) > 0 {
		if linkData.WebUrl == "" {
			return fmt.Sprintf("%s&%s", config.AppConfig.DefaultFallbackUrl, urlData.Encode())
		}
		return fmt.Sprintf("%s&%s", linkData.WebUrl, urlData.Encode())
	}

	if linkData.WebUrl == "" {
		return fmt.Sprintf("%s?%s", config.AppConfig.DefaultFallbackUrl, urlData.Encode())
	}
	return fmt.Sprintf("%s?%s", linkData.WebUrl, urlData.Encode())
}

func frameAndroidUrl(android, shortCode string, utm map[string]string) string {
	var parsed MobileInputs
	json.Unmarshal([]byte(android), &parsed)
	if parsed.Fbl == "" || android == "" {
		if config.AppConfig.Android.Behaviour == APP_SEARCH {
			var play = config.AppConfig.Android.GooglePlaySearchUrl
			utm["shortcode"] = shortCode
			values := url.Values{}
			for key, value := range utm {
				values.Add(key, value)
			}
			return fmt.Sprintf("%s&referrer=%s", play, values.Encode())
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
