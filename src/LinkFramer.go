package main

import (
	"encoding/json"
)

func frameAndroidUrl(android string) string {
	var parsed MobileInputs
	json.Unmarshal([]byte(android), &parsed)
	if parsed.Fbl == "" {
		if config.AppConfig.Android.Behaviour == APP_SEARCH {
			return config.AppConfig.Android.GooglePlaySearchUrl
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
