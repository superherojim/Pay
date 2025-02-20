package task

import (
	"cheemshappy_pay/internal/service"
	"context"
	"fmt"
)

type OrderTask interface {
	CheckOrder(ctx context.Context) error
}

func NewOrderTask(
	task *Task,
	orderService service.OrderService,
) OrderTask {
	return &orderTask{
		Task:         task,
		orderService: orderService,
	}
}

type orderTask struct {
	orderService service.OrderService
	*Task
}

func (t orderTask) CheckOrder(ctx context.Context) error {
	// do something
	fmt.Println("CheckOrder")
	t.orderService.ListenOrderPay(ctx)
	return nil
}
