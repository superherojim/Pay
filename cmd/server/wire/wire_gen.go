// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"cheemshappy_pay/internal/handler"
	"cheemshappy_pay/internal/repository"
	"cheemshappy_pay/internal/server"
	"cheemshappy_pay/internal/service"
	"cheemshappy_pay/pkg/app"
	"cheemshappy_pay/pkg/chain"
	"cheemshappy_pay/pkg/helper/sid"
	"cheemshappy_pay/pkg/jwt"
	"cheemshappy_pay/pkg/log"
	"cheemshappy_pay/pkg/server/http"
	"github.com/google/wire"
	"github.com/imroc/req/v3"
	"github.com/spf13/viper"
)

// Injectors from wire.go:

func NewWire(viperViper *viper.Viper, logger *log.Logger) (*app.App, func(), error) {
	client := repository.NewRedis(viperViper)
	jwtJWT := jwt.NewJwt(viperViper, client)
	handlerHandler := handler.NewHandler(logger)
	db := repository.NewDB(viperViper, logger)
	repositoryRepository := repository.NewRepository(db, client, logger)
	transaction := repository.NewTransaction(repositoryRepository)
	sidSid := sid.NewSid()
	serviceService := service.NewService(transaction, logger, sidSid, jwtJWT)
	reqClient := req.NewClient()
	userRepository := repository.NewUserRepository(repositoryRepository, reqClient)
	userService := service.NewUserService(serviceService, userRepository, reqClient)
	userHandler := handler.NewUserHandler(handlerHandler, userService)
	walletRepository := repository.NewWalletRepository(repositoryRepository)
	sysWalletRepository := repository.NewSysWalletRepository(repositoryRepository)
	walletService := service.NewWalletService(walletRepository, sysWalletRepository)
	walletHandler := handler.NewWalletHandler(handlerHandler, walletService)
	merchantsRepository := repository.NewMerchantsRepository(repositoryRepository)
	merchantsMetaRepository := repository.NewMerchantsMetaRepository(repositoryRepository)
	merchantsMetaService := service.NewMerchantsMetaService(serviceService, merchantsMetaRepository)
	merchantsService := service.NewMerchantsService(serviceService, merchantsRepository, merchantsMetaService, walletService)
	merchantsHandler := handler.NewMerchantsHandler(handlerHandler, merchantsService)
	merchantsMetaHandler := handler.NewMerchantsMetaHandler(handlerHandler, merchantsMetaService)
	orderRepository := repository.NewOrderRepository(repositoryRepository, client)
	merchantsApiRepository := repository.NewMerchantsApiRepository(repositoryRepository, client)
	merchantsApiService := service.NewMerchantsApiService(serviceService, merchantsApiRepository)
	sysConfigRepository := repository.NewSysConfigRepository(repositoryRepository)
	sysConfigService := service.NewSysConfigService(serviceService, sysConfigRepository)
	verifierFactory := chain.NewVerifierFactory()
	orderService := service.NewOrderService(serviceService, orderRepository, merchantsApiService, merchantsService, walletService, sysConfigService, viperViper, verifierFactory)
	orderHandler := handler.NewOrderHandler(handlerHandler, orderService, walletService, merchantsApiService)
	sysWalletService := service.NewSysWalletService(sysWalletRepository)
	sysWalletHandler := handler.NewSysWalletHandler(handlerHandler, sysWalletService, walletService)
	statsService := service.NewStatsService(merchantsRepository, orderRepository)
	statsHandler := handler.NewStatsHandler(statsService)
	merchantsApiHandler := handler.NewMerchantsApiHandler(handlerHandler, merchantsApiService)
	sysConfigHandler := handler.NewSysConfigHandler(handlerHandler, sysConfigService)
	httpServer := server.NewHTTPServer(logger, viperViper, jwtJWT, userHandler, walletHandler, merchantsHandler, merchantsMetaHandler, orderHandler, sysWalletHandler, statsHandler, merchantsApiHandler, sysConfigHandler)
	job := server.NewJob(logger)
	task := server.NewTask(logger)
	appApp := newApp(httpServer, job, task)
	return appApp, func() {
	}, nil
}

// wire.go:

var repositorySet = wire.NewSet(repository.NewDB, repository.NewRedis, repository.NewRepository, repository.NewTransaction, repository.NewUserRepository, repository.NewWalletRepository, repository.NewMerchantsRepository, repository.NewMerchantsMetaRepository, repository.NewOrderRepository, repository.NewMerchantsApiRepository, repository.NewSysWalletRepository, repository.NewSysConfigRepository)

var serviceSet = wire.NewSet(service.NewService, service.NewUserService, service.NewWalletService, service.NewMerchantsApiService, service.NewMerchantsMetaService, service.NewMerchantsService, service.NewOrderService, service.NewSysWalletService, service.NewStatsService, service.NewSysConfigService)

var handlerSet = wire.NewSet(handler.NewHandler, handler.NewUserHandler, handler.NewWalletHandler, handler.NewMerchantsHandler, handler.NewMerchantsMetaHandler, handler.NewOrderHandler, handler.NewMerchantsApiHandler, handler.NewSysWalletHandler, handler.NewStatsHandler, handler.NewSysConfigHandler)

var serverSet = wire.NewSet(server.NewHTTPServer, server.NewJob, server.NewTask)

// build App
func newApp(httpServer *http.Server, job *server.Job, task *server.Task) *app.App {
	return app.NewApp(app.WithServer(httpServer, job, task), app.WithName("demo-server"))
}
