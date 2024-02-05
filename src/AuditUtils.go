package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func (a *ServiceAudit) SetBuild(sBuild string, usrAgent string) {
	b := Build{}
	b.UserAgent = usrAgent
	if len(sBuild) != 0 {
		err := json.Unmarshal([]byte(sBuild), &b)
		if err != nil {
			WarningLogger.Printf("Invalid JSON; build %v", err)
		}
	}
	a.Build = b
}

func (audit *ServiceAudit) setRespDataToAudit(resp *Resp) {
	if !audit.SkipResponse {
		audit.Res = resp
	} else {
		audit.Res = SKIP_RESPONSE
	}
}

func (audit *ServiceAudit) SetSC(sc int) {
	audit.Sc = strconv.Itoa(sc)
}

func (audit *ServiceAudit) setPostData(rq *http.Request) {
	var err error
	var bodyBytes []byte
	if IsFormRequest(rq) {
		var reqFields = make(map[string]string)
		err = rq.ParseForm()
		for key := range rq.PostForm {
			reqFields[key] = rq.PostForm.Get(key)
		}
		bodyBytes, err = json.Marshal(reqFields)
	} else if IsJsonRequest(rq) {
		bodyBytes, _ = ioutil.ReadAll(rq.Body)
		rq.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	if err == nil {
		audit.Fields = json.RawMessage(bodyBytes)
	}
}

func (audit *ServiceAudit) SetFields(rq *http.Request) {
	if rq.Method == http.MethodGet || rq.Method == http.MethodDelete {
		var reqFields = make(map[string]string)
		for key := range rq.URL.Query() {
			reqFields[key] = rq.URL.Query().Get(key)
		}
		jsonVal, err := json.Marshal(reqFields)
		if err == nil {
			audit.Fields = json.RawMessage(jsonVal)
		}
	} else {
		audit.setPostData(rq)
	}
}

func (audit *ServiceAudit) SkipResponseAudit() {
	audit.SkipResponse = true
}

func (audit *ServiceAudit) SetErrMsgToContext(errorMsg string, ty string) {
	audit.Err = LogException{Type: ty, Msg: errorMsg, Status: FALSE}
}

func (audit *ServiceAudit) AppendErrListToContext(err string) {
	if audit.ErrList == nil {
		audit.ErrList = make([]string, 0)
	}
	audit.ErrList = append(audit.ErrList, err)
}

func (audit *ServiceAudit) LogAudit(startTime time.Time) {
	audit.ResT = CurrentTime()
	audit.TT = time.Since(startTime).Milliseconds()

	audit.SvSt = TRUE
	if audit.Err != (LogException{}) {
		audit.SvSt = FALSE
	}

	if len(audit.ErrList) > 0 {
		audit.SvSt = FALSE
	}

	jsonStr, _ := json.Marshal(audit)

	InfoLogger.Printf("AUDIT=%s", jsonStr)
}

func (audit *ServiceAudit) SetMsg(msg string) {
	audit.Msg = msg
}

func (audit *ServiceAudit) MarkAudit() {
	audit.LogAudit(GetDtInTime(STD_TIME_FORMAT, audit.ReqT))
}
