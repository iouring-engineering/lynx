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

// func (db *LynxDbType) RetryDbConn() {
// 	go func() {
// 		ticker := time.NewTicker(5 * time.Second)
// 		for range ticker.C {
// 			err := LynxDb.Db.Ping()
// 			if err != nil {
// 				InfoLogger.Println("Lost connection to the database. Reconnecting...")
// 				err = db.InitLynxDbConn()
// 				if err != nil {
// 					InfoLogger.Println("Failed to reconnect:", err)
// 				} else {
// 					InfoLogger.Println("Reconnected to the database.")
// 				}
// 			}
// 		}
// 	}()
// }

func (db *LynxDbType) InitLynxDbConn() error {
	const dbDriver = "postgres"
	dbUser := config.LynxDb.User
	dbPass := config.LynxDb.Password
	dbName := config.LynxDb.Database
	dbHost := config.LynxDb.Host
	dbPort := config.LynxDb.Port
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)
	var err error
	localDb, err := sql.Open(dbDriver, connStr)

	if err != nil {
		ErrorLogger.Panicln(err.Error())
		return err
	}
	InfoLogger.Println("DB configuration setup done")
	localDb.SetMaxIdleConns(config.LynxDb.MaxIdle)
	localDb.SetMaxOpenConns(config.LynxDb.MaxOpen)

	LynxDb = &LynxDbType{Db: localDb, Endpoint: "MKT-DATA"}
	return nil
}

func (db *LynxDbType) getEndPoint(cxt *IouHttpContext) EndPointContext {
	return EndPointContext{Db: LynxDb.Db, EndpointName: LynxDb.Endpoint, IouHttpContext: cxt}
}

// func (db *LynxDbType) checkConn(cxt *IouHttpContext) error {
// 	end := db.getEndPoint(cxt)
// 	err := end.Db.Ping()
// 	return err
// }

func (db *LynxDbType) insertShortLink(cxt *IouHttpContext, input InsertShortLink) error {
	_, err := IODBPrepareExec(db.getEndPoint(cxt), INSERT_SHORT_LINK, input.ShortCode, input.Data, input.WebUrl,
		input.Android, input.Ios, input.Desktop, input.Social, input.Expiry)
	return err
}
