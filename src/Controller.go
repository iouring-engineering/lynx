package main

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
	err := ReqFormDecode(cxt, &request)
	if err != nil {

	}
	if err = validateCreateLink(request); err != nil {
		cxt.SendBadReqResponse(err.Error())
		return
	}
	var res = &Resp{
		S:   RESP_OK,
		Msg: "ok",
	}
	cxt.SendResponse(res)
}
