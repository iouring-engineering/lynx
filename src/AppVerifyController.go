package main

import "fmt"

func IosVerify(cxt *IouHttpContext) {
	var resp IosAppVerifyResponse
	var appLinks IosAppLinks
	appLinks.Apps = make([]string, 0)
	appLinks.Details = make([]IosAppDetails, 0)
	var appDetails IosAppDetails
	appDetails.AppId = fmt.Sprintf("%s.%s", config.AppConfig.Ios.TeamId, config.AppConfig.Ios.BundleIdentifier)
	appDetails.Paths = config.AppConfig.Ios.AppLinkPath
	appLinks.Details = append(appLinks.Details, appDetails)
	resp.AppLinks = appLinks
	cxt.SendAnyResponse(resp)
}

func AndroidVerify(cxt *IouHttpContext) {
	var resp []AndroidVerifyResponse = make([]AndroidVerifyResponse, 0)
	var respObj AndroidVerifyResponse
	var target AndroidTarget
	target.NameSpace = ANDROID_NAMESPACE
	target.PackageName = config.AppConfig.Android.PackageName
	target.Sha256 = config.AppConfig.Android.Certificate
	respObj.Relation = []string{ANDROID_RELATION}
	respObj.Target = target
	resp = append(resp, respObj)
	cxt.SendAnyResponse(resp)
}
