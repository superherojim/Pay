package task

import (
	"cheemshappy_pay/internal/service"
	"context"
)

type OrderTask interface {
	CheckListenOrder(ctx context.Context) error
	CheckPendingOrder(ctx context.Context) error
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

func (t orderTask) CheckListenOrder(ctx context.Context) error {
	t.orderService.ListenOrderPay(ctx)
	return nil
}

func (t orderTask) CheckPendingOrder(ctx context.Context) error {
	t.orderService.CheckPendingOrder(ctx)
	return nil
}
