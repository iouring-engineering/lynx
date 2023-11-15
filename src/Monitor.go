package main

import (
	"encoding/json"
)

// var monitorLock sync.Mutex

var (
	endpointMap = make(map[string]bool)
)

type MonitorEx struct {
	Msg string `json:"msg,omitempty"`
	Ty  string `json:"ty,omitempty"`
}

type MonitorData struct {
	SvSt bool      `json:"svSt"`
	Msg  string    `json:"msg,omitempty"`
	Ty   string    `json:"ty"`
	Time string    `json:"resT"`
	Ex   MonitorEx `json:"ex,omitempty"`
}

func MarkSuccess(endCnxt EndPointContext) {
	if endCnxt.LogOnce {
		val, exists := endpointMap[endCnxt.EndpointName]
		// already success
		if !exists || val {
			return
		}
		endpointMap[endCnxt.EndpointName] = true
		mData := &MonitorData{
			Ty:   endCnxt.EndpointName,
			Time: CurrentTime(),
			SvSt: true,
		}
		logMonitorData(mData)
	}
}

func MarkFailure(endCnxt EndPointContext, msg string) {
	if endCnxt.LogOnce {
		val, exists := endpointMap[endCnxt.EndpointName]

		if exists && !val {
			return
		}
		endpointMap[endCnxt.EndpointName] = false
		mEx := MonitorEx{
			Msg: msg,
			Ty:  endCnxt.EndpointName,
		}

		mData := &MonitorData{
			Ty:   endCnxt.EndpointName,
			Time: CurrentTime(),
			SvSt: false,
			Ex:   mEx,
		}
		logMonitorData(mData)
	} else {
		InfoLogger.Println(endCnxt.EndpointName, msg)
	}
}

func logMonitorData(mData *MonitorData) {
	jsonResp, err := json.Marshal(mData)
	if err == nil {
		InfoLogger.Printf("AUDIT=%v", string(jsonResp))
	}
}
