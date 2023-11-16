package main

import (
	"github.com/gorilla/mux"
)

// @title						Lynx Services
// @version						0.0.1
// @description					Micro Module to create and share short links
// @tag.name					Links
// @tag.description				Creating and sharing short links
// @host 						localhost:8080
// @schemes 					http
// @BasePath					/lynx
// @securityDefinitions.apikey	http_bearer
// @in 							header
// @name 						Authorization
func main() {
	config = &Config{}
	InitializeConfigs(config)
	InfoLogger.Println("configs initialized")
	LynxDb.InitLynxDbConn()
	InfoLogger.Println("Lynx DB connected")
	err := loadHtmlFile()
	if err != nil {
		InfoLogger.Println("error on html ", err)
		return
	}
	var muxRouter *mux.Router = mux.NewRouter()
	var localRouter = &Router{}
	localRouter.initializeRouter(muxRouter)
	InitializeHttpServer(muxRouter)
}
