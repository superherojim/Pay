package server

import (
	"cheemshappy_pay/internal/task"
	"context"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
)

type TaskServer struct {
	scheduler *gocron.Scheduler
	orderTask task.OrderTask
}

func NewTaskServer(
	orderTask task.OrderTask,
) *TaskServer {
	return &TaskServer{
		orderTask: orderTask,
	}
}
func (t *TaskServer) Start(ctx context.Context) error {
	gocron.SetPanicHandler(func(jobName string, recoverData interface{}) {
		fmt.Println("TaskServer Panic", zap.String("job", jobName), zap.Any("recover", recoverData))
	})

	// eg: crontab task
	t.scheduler = gocron.NewScheduler(time.UTC)
	// if you are in China, you will need to change the time zone as follows
	// t.scheduler = gocron.NewScheduler(time.FixedZone("PRC", 8*60*60))

	//_, err := t.scheduler.Every("3s").Do(func()
	_, err := t.scheduler.Every("60s").Do(func() {
		err := t.orderTask.CheckOrder(ctx)
		if err != nil {
			fmt.Println("CheckOrder error", zap.Error(err))
		}
	})
	if err != nil {
		fmt.Println("CheckOrder error", zap.Error(err))
	}

	t.scheduler.StartBlocking()
	return nil
}
func (t *TaskServer) Stop(ctx context.Context) error {
	t.scheduler.Stop()
	fmt.Println("TaskServer stop...")
	return nil
}
