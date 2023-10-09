package main

import (
	"encoding/json"

	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

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
