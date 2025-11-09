package main

import (
	"github.com/shakuni-dyutas/dyutas-auth/internal/app/adpt/googleadpt"
	"github.com/shakuni-dyutas/dyutas-auth/internal/app/repo/userrepo"
	"github.com/shakuni-dyutas/dyutas-auth/internal/app/svc/authsvc"
	"github.com/shakuni-dyutas/dyutas-auth/internal/config"
	"github.com/shakuni-dyutas/dyutas-auth/internal/db"
	"github.com/shakuni-dyutas/dyutas-auth/internal/log"
	"github.com/shakuni-dyutas/dyutas-auth/internal/rest"
)

func main() {
	logConfig := config.LoadLogConfig()
	logger := log.NewLogger(logConfig)

	appConfig, err := config.LoadAppConfigs()
	if err != nil {
		logger.Error("failed to load app configs", "error", err)
		return
	}

	dbConn, err := db.NewRDBConnectionPool(db.ConnectionConfig{
		User: appConfig.RDbConfig.User,
		Pw:   appConfig.RDbConfig.Pw,
		Host: appConfig.RDbConfig.Host,
		Port: appConfig.RDbConfig.Port,
		DB:   appConfig.RDbConfig.Db,
	})
	if err != nil {
		logger.Error("failed to create RDB connection pool", "error", err)
		return
	}

	googleAdptConfig := googleadpt.GoogleAdptConfig{
		ClientId:     appConfig.AuthConfig.GoogleClientId,
		ClientSecret: appConfig.AuthConfig.GoogleClientSecret,
	}
	googleAdpt := googleadpt.New(googleAdptConfig)

	userRepo := userrepo.New(dbConn)
	authSvc := authsvc.New(googleAdpt, userRepo, appConfig.AuthConfig.AppJwtKey)

	// with go routine run rest server and exit if err is returned
	restErrChan := make(chan error)

	go func() {
		err = rest.Run(rest.RestConfig{
			Port:          appConfig.Port,
			AllowOrigins:  appConfig.AllowOrigins,
			AllowMethods:  appConfig.AllowMethods,
			AllowHeaders:  appConfig.AllowHeaders,
			ExposeHeaders: appConfig.ExposeHeaders,
			Logger:        logger,
		}, authSvc)
		if err != nil {
			restErrChan <- err
		}
	}()

	if err := <-restErrChan; err != nil {
		logger.Error("REST server failed", "error", err)
		return
	}
}
