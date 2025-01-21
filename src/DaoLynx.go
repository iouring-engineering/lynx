package main

import (
	"database/sql"
	"fmt"
)

type LynxDbType DatabaseType

var (
	LynxDb *LynxDbType
)

func (db *LynxDbType) InitLynxDbConn() error {
	const dbDriver = "mysql"
	dbUser := config.LynxDb.User
	dbPass := config.LynxDb.Password
	dbName := config.LynxDb.Database
	dbHost := config.LynxDb.Host
	dbPort := config.LynxDb.Port
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	var err error
	localDb, err := sql.Open(dbDriver, connStr)

	if err != nil {
		Logger.Error(err.Error())
		return err
	}
	Logger.Info("DB configuration setup done")
	localDb.SetMaxIdleConns(config.LynxDb.MaxIdle)
	localDb.SetMaxOpenConns(config.LynxDb.MaxOpen)

	LynxDb = &LynxDbType{Db: localDb, Endpoint: "MKT-DATA"}
	return nil
}

func (db *LynxDbType) getEndPoint(cxt *IouHttpContext) EndPointContext {
	return EndPointContext{Db: LynxDb.Db, EndpointName: LynxDb.Endpoint, IouHttpContext: cxt}
}

func (db *LynxDbType) insertShortLink(cxt *IouHttpContext, input DbShortLink) error {
	_, err := IODBPrepareExec(db.getEndPoint(cxt), INSERT_SHORT_LINK, input.ShortCode, input.Data, input.WebUrl,
		input.Android, input.Ios, input.Social, input.Expiry)
	return err
}

func (db *LynxDbType) getData(cxt *IouHttpContext, shortCode string) (DbShortLink, bool, error) {
	var res DbShortLink
	rows, err := IODBPrepareQuery(db.getEndPoint(cxt), GET_LINK_DATA, shortCode)
	if err != nil {
		return res, false, err
	}
	defer rows.Close()
	if rows.Next() {
		res.ShortCode = shortCode
		err = rows.Scan(&res.Data, &res.WebUrl, &res.Android, &res.Ios, &res.Social, &res.Expiry)
		if err != nil {
			return res, false, err
		}
		return res, true, nil
	}
	return res, false, nil
}
