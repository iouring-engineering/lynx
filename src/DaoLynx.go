package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type LynxDbType DatabaseType

var (
	LynxDb *LynxDbType
)

func (db *LynxDbType) InitLynxDbConn() {
	const dbDriver = "postgres"
	dbUser := config.LynxDb.User
	dbPass := config.LynxDb.Password
	dbName := config.LynxDb.Database
	dbHost := config.LynxDb.Host
	dbPort := config.LynxDb.Port
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)
	fmt.Println(connStr)
	var err error
	localDb, err := sql.Open(dbDriver, connStr)
	localDb.SetMaxIdleConns(config.LynxDb.MaxIdle)
	localDb.SetMaxOpenConns(config.LynxDb.MaxOpen)

	LynxDb = &LynxDbType{Db: localDb, Endpoint: "MKT-DATA"}

	InfoLogger.Println("DB configuration setup done")
	if err != nil {
		ErrorLogger.Panicln(err.Error())
	}
}

func (db *LynxDbType) getEndPoint(cxt *IouHttpContext) EndPointContext {
	return EndPointContext{Db: LynxDb.Db, EndpointName: LynxDb.Endpoint, IouHttpContext: cxt}
}

func (db *LynxDbType) checkConn(cxt *IouHttpContext) error {
	end := db.getEndPoint(cxt)
	err := end.Db.Ping()
	return err
}
