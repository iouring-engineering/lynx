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

	}
	if err = validateCreateLink(request); err != nil {
		cxt.SendBadReqResponse(err.Error())
		return
	}
	androidStr, _ := json.Marshal(request.Android)
	iosStr, _ := json.Marshal(request.Ios)
	socialStr, _ := json.Marshal(request.Social)
	var expiryValue = calculateExpiry(request.Expiry)
	var shortCode string
	var retryCount = 0
	strData, err := getDataString(request.Data)
	if err != nil {
		cxt.SendBadReqResponse("Invalid data")
		return
	}
	for {
		shortCode = genShortUrl()
		var dbExp = sql.NullString{String: "", Valid: false}
		if expiryValue != 0 {

		}
		err = LynxDb.insertShortLink(cxt, DbShortLink{
			ShortCode: shortCode,
			Data:      strData,
			WebUrl:    request.WebUrl,
			Android:   string(androidStr),
			Ios:       string(iosStr),
			Social:    string(socialStr),
			Expiry:    dbExp,
		})
		retryCount++
		if err != nil && isDuplicateLink(err) {
			if retryCount <= config.AppConfig.DuplicateRetryCount {
				continue
			} else {
				var res = &Resp{
					S:   RESP_ERROR,
					Msg: ShortLinkFailed,
				}
				cxt.SendResponse(res)
				return
			}
		}
		break
	}

	var res = &Resp{
		S: RESP_OK,
		D: map[string]string{
			"link": fmt.Sprintf("%s/%s", config.AppConfig.BaseUrl, shortCode),
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
	linkData, exists, err := LynxDb.getData(cxt, shortCode)
	if err != nil {
		return
	}
	if !exists {
		html := frame404WebPage()
		cxt.sendHtmlResponse(html)
		return
	}
	if isDesktopWeb(cxt) {
		var url string = frameBrowserUrl(linkData)
		cxt.SendRedirect(url)
		return
	}

	if isAndroidWeb(cxt) {
		var url string = frameAndroidUrl(linkData.Android)
		cxt.SendRedirect(url)
		return
	}

	if isIosWeb(cxt) {
		var url string = frameIosUrl(linkData.Ios)
		cxt.SendRedirect(url)
		return
	}

	html := frameWebPage(linkData, config.AppConfig.DefaultFallbackUrl)
	cxt.sendHtmlResponse(html)
}

// @Summary		Get data
// @Description	Get data using short code
// @Tags        Links
// @Id 			get-source-link-data
// @Success		200 {object} string
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
		return
	}
	if !exists {
		cxt.SendNoDataResponse()
		return
	}
	if isValidJson(linkData.Data) {
		var res *Resp = &Resp{S: RESP_OK, D: json.RawMessage(linkData.Data)}
		cxt.SendResponse(res)
		return
	}
	var res *Resp = &Resp{S: RESP_OK, D: linkData.Data}
	cxt.SendResponse(res)
}
