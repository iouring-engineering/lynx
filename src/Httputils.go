package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS,PUT")

		w.Header().Set("Access-Control-Allow-Headers", "accept, If-None-Match, Origin, Content-Type, "+
			"X-TOKENDATA, X-BUILD, x-tokendata,cache-control, content-type, DNT, Keep-Alive, User-Agent,"+
			" X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Set-Cookie, origin, accept, Authorization")

		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func ReqBodyDecode[S any](cxt *IouHttpContext, reqModel *S) error {
	err := json.NewDecoder(cxt.Request.Body).Decode(reqModel)
	return err
}
func ReqFormDecode[S any](cxt *IouHttpContext, reqModel *S) error {
	err := cxt.Request.ParseForm()
	if err != nil {
		return err
	}
	decoder.IgnoreUnknownKeys(true)
	err = decoder.Decode(reqModel, cxt.Request.PostForm)
	if err != nil {
		return err
	}
	return nil
}

func (cxt *IouHttpContext) SendResponse(resp *Resp) {
	jsonResp, err := JSONMarshal(resp)

	if err != nil {
		Logger.Error(fmt.Printf("Error happened in JSON marshal. Err: %v", err))
	}
	cxt.Audit.setRespDataToAudit(resp)

	cxt.RespWriter.Header().Set("Content-Type", "application/json")
	cxt.RespWriter.Write(jsonResp)
}

func (cxt *IouHttpContext) SendAnyResponse(resp any) {
	jsonResp, err := JSONMarshal(resp)

	if err != nil {
		Logger.Error(fmt.Printf("Error happened in JSON marshal. Err: %v", err))
	}

	cxt.RespWriter.Header().Set("Content-Type", "application/json")
	cxt.RespWriter.Write(jsonResp)
}

func (cxt *IouHttpContext) SendRedirect(url string) {
	http.Redirect(cxt.RespWriter, cxt.Request, url, http.StatusMovedPermanently)
}

func (cxt *IouHttpContext) SendErrResponse(httpStatusCode int, message string) {
	var resp = &Resp{S: "error", Msg: message}
	jsonResp, err := json.Marshal(resp)

	if err != nil {
		Logger.Error(fmt.Sprintf("Error happened in JSON marshal. Err: %v", err))
	}

	cxt.RespWriter.Header().Set("Content-Type", APPLICATION_JSON)
	cxt.RespWriter.WriteHeader(httpStatusCode)
	cxt.RespWriter.Write(jsonResp)
}

func (cxt *IouHttpContext) SendBadReqResponse(message string) {
	cxt.SendErrResponse(http.StatusBadRequest, message)
}

func (cxt *IouHttpContext) SendUnAuthResponse() {
	cxt.SendErrResponse(http.StatusUnauthorized, UnAuthorized)
}

func (cxt *IouHttpContext) SendNoDataResponse() {
	var NoDataResp *Resp = &Resp{
		S:   RESP_NODATA,
		Msg: "No Data available.",
	}
	cxt.SendResponse(NoDataResp)
}

func (cxt *IouHttpContext) SendResponseMsg(status string, msg string) {
	resp := &Resp{
		S:   status,
		Msg: msg,
	}

	cxt.SendResponse(resp)
}

func (cxt *IouHttpContext) sendHtmlResponse(respBytes string) {
	cxt.Audit.Res = respBytes
	cxt.RespWriter.Header().Set("Content-Type", "text/html")
	cxt.RespWriter.Write([]byte(respBytes))
}

func InitializeHttpServer(router *mux.Router) {
	if libConfig.Env == LOCAL {
		router.Use(CorsMiddleware)
		router.PathPrefix("/").Handler(http.FileServer(http.Dir("./.service-validations")))
	}

	var server *http.Server = &http.Server{
		Handler:      router,
		Addr:         CmdArgs.ServerAddr,
		WriteTimeout: time.Duration(libConfig.Http.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(libConfig.Http.ReadTimeout) * time.Second,
	}
	fmt.Printf("Http Server started and listening on %s\n", CmdArgs.ServerAddr)
	server.ListenAndServe()
}
