package main

import (
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
	desktopStr, _ := json.Marshal(request.Desktop)
	socialStr, _ := json.Marshal(request.Social)
	var expiryValue = calculateExpiry(string(request.Expiry.Type), request.Expiry.Value)
	var shortCode string
	var retryCount = 0
	strData, err := getDataString(request.Data)
	if err != nil {
		cxt.SendBadReqResponse("Invalid data")
		return
	}
	for {
		shortCode = genShortUrl()
		err = LynxDb.insertShortLink(cxt, DbShortLink{
			ShortCode: shortCode,
			Data:      strData,
			WebUrl:    request.WebUrl,
			Android:   string(androidStr),
			Ios:       string(iosStr),
			Desktop:   string(desktopStr),
			Social:    string(socialStr),
			Expiry:    expiryValue,
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
// @Param       shorturl   path  string true "shorturl" example(CJloO)
// @Router      /{shorturl} [get]
func GetSourceLink(cxt *IouHttpContext) {
	req := mux.Vars(cxt.Request)
	shortCode := req["shorturl"]
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
		sendHtmlResponse(cxt, []byte(html))
		return
	}
	if isDesktopWeb(cxt) {
		var url string = frameDesktopBrowser(linkData)
		html := frameWebPage(linkData, url)
		sendHtmlResponse(cxt, []byte(html))
		return
	}
}
