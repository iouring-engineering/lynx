package main

import "encoding/json"

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
	LynxDb.insertShortLink(cxt, InsertShortLink{
		ShortCode: genShortUrl(),
		Data:      request.Data,
		WebUrl:    request.WebUrl,
		Android:   string(androidStr),
		Ios:       string(iosStr),
		Desktop:   string(desktopStr),
		Social:    string(socialStr),
		Expiry:    expiryValue,
	})
	var res = &Resp{
		S:   RESP_OK,
		Msg: "ok",
	}
	cxt.SendResponse(res)
}
