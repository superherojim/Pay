//go:build wireinject
// +build wireinject

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

var repositorySet = wire.NewSet(
	repository.NewDB,
	repository.NewRedis,
	repository.NewRepository,
	repository.NewTransaction,
	repository.NewUserRepository,
	repository.NewWalletRepository,
	repository.NewMerchantsRepository,
	repository.NewMerchantsMetaRepository,
	repository.NewOrderRepository,
	repository.NewMerchantsApiRepository,
	repository.NewSysWalletRepository,
	repository.NewSysConfigRepository,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,
	service.NewWalletService,
	service.NewMerchantsApiService,
	service.NewMerchantsMetaService,
	service.NewMerchantsService,
	service.NewOrderService,
	service.NewSysWalletService,
	service.NewStatsService,
	service.NewSysConfigService,
)
var taskSet = wire.NewSet(
	task.NewTask,
	task.NewOrderTask,
)

var serverSet = wire.NewSet(
	server.NewTaskServer,
)

// build App
func newApp(
	task *server.TaskServer,
) *app.App {
	return app.NewApp(
		app.WithServer(task),
		app.WithName("demo-task"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		repositorySet,
		serviceSet,
		serverSet,
		sid.NewSid,
		jwt.NewJwt,
		newApp,
		chain.NewVerifierFactory,
		taskSet,
	))
}
