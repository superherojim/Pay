// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"cheemshappy_pay/internal/repository"
	"cheemshappy_pay/internal/server"
	"cheemshappy_pay/internal/service"
	"cheemshappy_pay/internal/task"
	"cheemshappy_pay/pkg/app"
	"cheemshappy_pay/pkg/chain"
	"cheemshappy_pay/pkg/helper/sid"
	"cheemshappy_pay/pkg/jwt"
	"cheemshappy_pay/pkg/log"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

// Injectors from wire.go:

func NewWire(viperViper *viper.Viper, logger *log.Logger) (*app.App, func(), error) {
	taskTask := task.NewTask()
	db := repository.NewDB(viperViper, logger)
	client := repository.NewRedis(viperViper)
	repositoryRepository := repository.NewRepository(db, client, logger)
	transaction := repository.NewTransaction(repositoryRepository)
	sidSid := sid.NewSid()
	jwtJWT := jwt.NewJwt(viperViper, client)
	serviceService := service.NewService(transaction, logger, sidSid, jwtJWT)
	orderRepository := repository.NewOrderRepository(repositoryRepository, client)
	merchantsApiRepository := repository.NewMerchantsApiRepository(repositoryRepository, client)
	merchantsApiService := service.NewMerchantsApiService(serviceService, merchantsApiRepository)
	merchantsRepository := repository.NewMerchantsRepository(repositoryRepository)
	merchantsMetaRepository := repository.NewMerchantsMetaRepository(repositoryRepository)
	merchantsMetaService := service.NewMerchantsMetaService(serviceService, merchantsMetaRepository)
	walletRepository := repository.NewWalletRepository(repositoryRepository)
	sysWalletRepository := repository.NewSysWalletRepository(repositoryRepository)
	walletService := service.NewWalletService(walletRepository, sysWalletRepository)
	merchantsService := service.NewMerchantsService(serviceService, merchantsRepository, merchantsMetaService, walletService)
	sysConfigRepository := repository.NewSysConfigRepository(repositoryRepository)
	sysConfigService := service.NewSysConfigService(serviceService, sysConfigRepository)
	verifierFactory := chain.NewVerifierFactory()
	orderService := service.NewOrderService(serviceService, orderRepository, merchantsApiService, merchantsService, walletService, sysConfigService, viperViper, verifierFactory)
	orderTask := task.NewOrderTask(taskTask, orderService)
	taskServer := server.NewTaskServer(orderTask)
	appApp := newApp(taskServer)
	return appApp, func() {
	}, nil
}

// wire.go:

var repositorySet = wire.NewSet(repository.NewDB, repository.NewRedis, repository.NewRepository, repository.NewTransaction, repository.NewUserRepository, repository.NewWalletRepository, repository.NewMerchantsRepository, repository.NewMerchantsMetaRepository, repository.NewOrderRepository, repository.NewMerchantsApiRepository, repository.NewSysWalletRepository, repository.NewSysConfigRepository)

var serviceSet = wire.NewSet(service.NewService, service.NewUserService, service.NewWalletService, service.NewMerchantsApiService, service.NewMerchantsMetaService, service.NewMerchantsService, service.NewOrderService, service.NewSysWalletService, service.NewStatsService, service.NewSysConfigService)

var taskSet = wire.NewSet(task.NewTask, task.NewOrderTask)

var serverSet = wire.NewSet(server.NewTaskServer)

// build App
func newApp(task2 *server.TaskServer,
) *app.App {
	return app.NewApp(app.WithServer(task2), app.WithName("demo-task"))
}
