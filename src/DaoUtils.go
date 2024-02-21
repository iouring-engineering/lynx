package main

import (
	"database/sql"
	"errors"
)

func IODBExec(cnxt EndPointContext, query string, args ...any) (sql.Result, error) {
	if cnxt.Db == nil {
		MarkFailure(cnxt, "Connection failure")
		return nil, errors.New("connection failure")
	}
	r, err := cnxt.Db.Exec(query, args...)
	if err != nil {
		MarkFailure(cnxt, err.Error())
	} else {
		MarkSuccess(cnxt)
	}
	return r, err
}

func IODBPrepareExec(cnxt EndPointContext, query string, args ...any) (sql.Result, error) {
	if cnxt.Db == nil {
		MarkFailure(cnxt, "Connection failure")
		return nil, errors.New("connection failure")
	}
	stmt, err := cnxt.Db.Prepare(query)

	if err != nil {
		MarkFailure(cnxt, err.Error())
		return nil, err
	} else {
		MarkSuccess(cnxt)
	}

	defer stmt.Close()

	res, err := stmt.Exec(args...)
	if err != nil {
		MarkFailure(cnxt, err.Error())
		return nil, err
	} else {
		MarkSuccess(cnxt)
	}
	return res, err
}

func IODBPrepareQuery(cnxt EndPointContext, query string, args ...any) (*sql.Rows, error) {
	if cnxt.Db == nil {
		MarkFailure(cnxt, "Connection failure")
		return nil, errors.New("connection failure")
	}
	stmt, err := cnxt.Db.Prepare(query)
	if err != nil {
		MarkFailure(cnxt, err.Error())
		return nil, err
	} else {
		MarkSuccess(cnxt)
	}

	defer stmt.Close()

	res, err := stmt.Query(args...)
	if err != nil {
		MarkFailure(cnxt, err.Error())
		return nil, err
	} else {
		MarkSuccess(cnxt)
	}
	return res, err
}
