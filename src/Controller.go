package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// @Summary		Create
// @Description	Create short link urls
// @Tags        Links
// @Id 			create-short-link
// @Accept      json
// @Success		200  {object} CreateShortLinkResponse
// @Produce     json
// @Param request body CreateShortLinkRequest true "Request body"
// @Security 	http_bearer
// @Router      /create [post]
func CreateShortLink(cxt *IouHttpContext) {
	var request CreateShortLinkRequest
	err := ReqBodyDecode(cxt, &request)
	if err != nil {
		cxt.SendBadReqResponse(err.Error())
		return
	}
	if err = validateCreateLink(request); err != nil {
		cxt.SendBadReqResponse(err.Error())
		return
	}
	androidStr, _ := json.Marshal(request.Android)
	iosStr, _ := json.Marshal(request.Ios)
	socialStr, _ := json.Marshal(request.Social)
	if !validateExpiry(&request.Expiry) {
		cxt.SendBadReqResponse("Invalid expiry")
		return
	}
	var shortCode string
	var retryCount = 0
	strData, err := getDataString(request.Data)
	if err != nil {
		cxt.SendBadReqResponse("Invalid data")
		return
	}

	for {
		shortCode = genShortUrl()
		var dbExp = sql.NullString{Valid: false}
		if request.Expiry != "" {
			dbExp = sql.NullString{String: request.Expiry, Valid: true}
		}
		err = LynxDb.insertShortLink(cxt, DbShortLink{Data: strData, WebUrl: request.WebUrl,
			Android: string(androidStr), Ios: string(iosStr), Social: string(socialStr),
			Expiry: dbExp, ShortCode: shortCode})
		retryCount++
		if err != nil && isDuplicateLink(err) {
			if retryCount <= config.AppConfig.DuplicateRetryCount {
				continue
			} else {
				var res = &Resp{S: RESP_ERROR, Msg: ShortLinkFailed}
				cxt.SendResponse(res)
				return
			}
		}
		if err != nil {
			cxt.SendErrResponse(http.StatusOK, EndpointErr)
			return
		}
		break
	}

	var res = &Resp{
		S: RESP_OK,
		D: map[string]string{
			"link":      fmt.Sprintf("%s/%s", config.AppConfig.BaseUrl, shortCode),
			"shortcode": shortCode,
		},
	}
	cxt.SendResponse(res)
}

// @Summary		Get Source url
// @Description	Get actual Source url with data
// @Tags        Links
// @Id 			get-source-link
// @Success		302
// @Success		200 {object} IouMsgResp
// @Produce     html
// // @Param        X-TOKENDATA header string true "Auth data"
// @Param       shortcode   path  string true "shortcode" example(CJloO)
// @Router      /{shortcode} [get]
func GetSourceLink(cxt *IouHttpContext) {
	req := mux.Vars(cxt.Request)
	shortCode := req["shortcode"]
	if shortCode == "" {
		cxt.SendErrResponse(http.StatusOK, InvalidShortUrl)
		return
	}
	linkData, _, err := LynxDb.getData(cxt, shortCode)
	if err != nil {
		cxt.SendErrResponse(http.StatusOK, InvalidShortUrl)
		return
	}

	if isAndroidWeb(cxt) {
		var url string = frameAndroidUrl(linkData.Android, shortCode)
		html := frameAndroidWebPage(linkData, url)
		cxt.sendHtmlResponse(html)
		return
	}

	if isIosWeb(cxt) {
		var url string = frameIosUrl(linkData.Ios)
		html := frameIosWebPage(linkData, url, shortCode)
		cxt.sendHtmlResponse(html)
		return
	}
	html := frameWebPage(linkData)
	cxt.sendHtmlResponse(html)
}

// @Summary		Get data
// @Description	Get data using short code
// @Tags        Links
// @Id 			get-source-link-data
// @Success		200 {object} ShortCodeDataResponse
// @Produce     json
// @Param       shortcode   path  string true "shortcode" example(CJloO)
// @Router      /data/{shortcode} [get]
func GetData(cxt *IouHttpContext) {
	req := mux.Vars(cxt.Request)
	shortCode := req["shortcode"]
	if shortCode == "" {
		cxt.SendErrResponse(http.StatusOK, InvalidShortUrl)
		return
	}
	linkData, exists, err := LynxDb.getData(cxt, shortCode)
	if err != nil {
		cxt.Audit.AppendErrListToContext(err.Error())
		cxt.SendNoDataResponse()
		return
	}
	if !exists {
		cxt.SendNoDataResponse()
		return
	}
	if isValidJson(linkData.Data) {
		var r = ShortCodeDataResponse{Input: json.RawMessage(linkData.Data), ShortCode: linkData.ShortCode}
		var res *Resp = &Resp{S: RESP_OK, D: r}
		cxt.SendResponse(res)
		return
	}
	var r = ShortCodeDataResponse{Input: linkData.Data, ShortCode: linkData.ShortCode}
	var res *Resp = &Resp{S: RESP_OK, D: r}
	cxt.SendResponse(res)
}
