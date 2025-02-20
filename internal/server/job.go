package server

import (
	"cheemshappy_pay/pkg/log"
	"context"
	"fmt"
)

type Job struct {
	log *log.Logger
}

func NewJob(log *log.Logger) *Job {
	return &Job{
		log: log,
	}
}

func (j *Job) Start(ctx context.Context) error {
	fmt.Println("timer init")
	//j.executor.RegTask("task.test", task.Test)
	//j.executor.RegTask("task.CancelOrder", j.orderService.CancelOrder)
	////j.executor.RegTask("task.deletePic", j.picService.DeletePicJob)
	//j.executor.RegTask("task.PreviewPic", j.fileService.DoPreView)
	//j.executor.Run()
	return nil
}
func (j *Job) Stop(ctx context.Context) error {
	//j.executor.Stop()
	return nil
}
