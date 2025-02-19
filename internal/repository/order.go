package repository

import (
	v1 "bk/api/v1"
	"bk/internal/model"
	"bk/pkg/enum"
	"bk/pkg/helper/xid"
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type OrderRepository interface {
	GetOrder(ctx context.Context, id int64) (*model.Order, error)
	GetOrderList(ctx context.Context, req *v1.OrderListReq) (*v1.Paginator, error)
	CreateOrder(ctx context.Context, or *model.Order, mid int64) (*model.Order, error)
	createOrderf(ctx context.Context, or *model.Order, mid int64, txhash ...string) (*model.Order, error)
	GetOrderNo() string
	CancelOrder(ctx context.Context, or *model.Order, mid int64) (*model.Order, error)
	cancelOrderf(ctx context.Context, or *model.Order, mid int64, txhash ...string) (*model.Order, error)
	SuccessOrder(ctx context.Context, or *model.Order, mid int64) (*model.Order, error)
	successOrderf(ctx context.Context, or *model.Order, mid int64, txhash ...string) (*model.Order, error)
	GetOrderLock(ctx context.Context, no string, mid int64) (bool, error)
	ReleaseOrderLock(ctx context.Context, no string) error
	ReTryOrderLock(ctx context.Context, or *model.Order, mid int64, fn func(ctx context.Context, oro *model.Order, mID int64, txhash ...string) (*model.Order, error), txhash ...string) error
	CountByStatusAndTime(ctx context.Context, status string, startTime time.Time) (int64, error)
	TotalCount(ctx context.Context) (int64, error)
	GetOrderPay(ctx context.Context, no string) (*model.Order, error)
	ListenOrder(ctx context.Context, no string, req *v1.OrderPayTxOut) (*model.Order, error)
	listenOrderf(ctx context.Context, or *model.Order, mid int64, txhash ...string) (*model.Order, error)
	UpdateOrder(ctx context.Context, or *model.Order) (*model.Order, error)
	updateOrderf(ctx context.Context, or *model.Order, mid int64, txhash ...string) (*model.Order, error)
	GetOrderByStatus(ctx context.Context, status string) ([]*model.Order, error)
}

func NewOrderRepository(
	repository *Repository,
	rdb *redis.Client,
) OrderRepository {
	return &orderRepository{
		Repository: repository,
		rdb:        rdb,
	}
}

type orderRepository struct {
	*Repository
	rdb *redis.Client
}

func (r *orderRepository) GetOrder(ctx context.Context, id int64) (*model.Order, error) {
	var order model.Order

	return &order, nil
}

func (r *orderRepository) GetOrderList(ctx context.Context, req *v1.OrderListReq) (*v1.Paginator, error) {
	tx := newOrder(r.DB(ctx))
	if req.OrderNo != "" {
		tx.Where(tx.No.Eq(req.OrderNo))
	}
	byPage, i, err := tx.FindByPage(req.GetOffset(), req.Size)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	p := &v1.Paginator{
		Total: i,
		Data:  byPage,
	}
	return p, nil
}
func (r *orderRepository) createOrderf(ctx context.Context, or *model.Order, mid int64, txhash ...string) (*model.Order, error) {
	tx := newOrder(r.DB(ctx))
	o, err := tx.Where(tx.No.Eq(or.No)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if o != nil {
		return nil, errors.New("order  already exists, please try again")
	}
	err = tx.Create(or)
	if err != nil {
		return nil, err
	}
	return or, nil
}

func (r *orderRepository) CreateOrder(ctx context.Context, or *model.Order, mid int64) (*model.Order, error) {
	or.No = r.GetOrderNo()
	err := r.ReTryOrderLock(ctx, or, mid, r.createOrderf)
	if err != nil {
		return nil, err
	}
	return or, nil
}

func (r *orderRepository) CancelOrder(ctx context.Context, or *model.Order, mid int64) (*model.Order, error) {
	err := r.ReTryOrderLock(ctx, or, mid, r.cancelOrderf)
	if err != nil {
		return nil, err
	}
	return or, nil
}

func (r *orderRepository) SuccessOrder(ctx context.Context, order *model.Order, mid int64) (*model.Order, error) {
	err := r.ReTryOrderLock(ctx, order, mid, r.successOrderf)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (r *orderRepository) cancelOrderf(ctx context.Context, or *model.Order, mid int64, txhash ...string) (*model.Order, error) {
	tx := newOrder(r.DB(ctx))
	o, err := tx.Where(tx.No.Eq(or.No), tx.Status.Eq(enum.OrderStatusPending), tx.MID.Eq(mid)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if o == nil {
		return nil, errors.New("order not found")
	}
	o.Status = enum.OrderStatusCanceled
	_, err = tx.Updates(o)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (r *orderRepository) successOrderf(ctx context.Context, or *model.Order, mid int64, txhash ...string) (*model.Order, error) {
	//TODO: 检查订单正确性
	tx := newOrder(r.DB(ctx))
	o, err := tx.Where(tx.No.Eq(or.No), tx.Status.Eq(enum.OrderStatusPending), tx.MID.Eq(mid)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if o == nil {
		return nil, errors.New("order not found")
	}
	o.Status = enum.OrderStatusSuccess
	_, err = tx.Updates(o)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (r *orderRepository) GetOrderNo() string {
	return xid.GenXID()
}

func (r *orderRepository) GetOrderLock(ctx context.Context, no string, mid int64) (bool, error) {
	return r.rdb.SetNX(ctx, no, true, time.Second*2).Result()
}

func (r *orderRepository) ReleaseOrderLock(ctx context.Context, no string) error {
	return r.rdb.Del(ctx, no).Err()
}

func (r *orderRepository) ReTryOrderLock(ctx context.Context, or *model.Order, mid int64, fn func(ctx context.Context, oro *model.Order, mID int64, txhash ...string) (*model.Order, error), txhash ...string) error {
	try := 0
	lock := true
	for !lock {
		try++
		if try > 4 {
			return errors.New("order already exists, please try again")
		}
		lock, err := r.GetOrderLock(ctx, or.No, mid)
		if err != nil {
			return err
		}
		if !lock {
			time.Sleep(time.Millisecond * 100)
		}
	}
	or, err := fn(ctx, or, mid, txhash...)
	if err != nil {
		return err
	}
	return nil
}

func (r *orderRepository) CountByStatusAndTime(ctx context.Context, status string, startTime time.Time) (int64, error) {
	tx := newOrder(r.DB(ctx))
	return tx.Where(tx.Status.Eq(status), tx.CreatedAt.Gt(startTime)).Count()
}

func (r *orderRepository) TotalCount(ctx context.Context) (int64, error) {
	tx := newOrder(r.DB(ctx))
	return tx.Count()
}

func (r *orderRepository) GetOrderPay(ctx context.Context, no string) (*model.Order, error) {
	tx := newOrder(r.DB(ctx))
	return tx.Where(tx.No.Eq(no), tx.Status.Eq(enum.OrderStatusPending)).First()
}

func (r *orderRepository) ListenOrder(ctx context.Context, no string, req *v1.OrderPayTxOut) (*model.Order, error) {
	tx := newOrder(r.DB(ctx))
	or, err := tx.Where(tx.CNo.Eq(no), tx.Status.Eq(enum.OrderStatusPending)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if or == nil {
		return nil, errors.New("order not found")
	}
	err = r.ReTryOrderLock(ctx, or, or.MID, r.listenOrderf, req.TxHash)
	if err != nil {
		return nil, err
	}
	return or, nil

}

func (r *orderRepository) listenOrderf(ctx context.Context, or *model.Order, mid int64, txhash ...string) (*model.Order, error) {
	if len(txhash) == 0 {
		return nil, errors.New("txhash is required")
	}
	tx := newOrder(r.DB(ctx))
	o, err := tx.Where(tx.No.Eq(or.No), tx.Status.Eq(enum.OrderStatusPending), tx.MID.Eq(mid)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if o == nil {
		return nil, errors.New("order not found")
	}
	o.TxHash = txhash[0]
	o.Status = enum.OrderStatusListening
	_, err = tx.Updates(o)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (r *orderRepository) UpdateOrder(ctx context.Context, or *model.Order) (*model.Order, error) {
	err := r.ReTryOrderLock(ctx, or, or.MID, r.updateOrderf)
	if err != nil {
		return nil, err
	}
	return or, nil
}

func (r *orderRepository) updateOrderf(ctx context.Context, or *model.Order, mid int64, txhash ...string) (*model.Order, error) {
	tx := newOrder(r.DB(ctx))
	o, err := tx.Where(tx.No.Eq(or.No), tx.MID.Eq(mid)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if o == nil {
		return nil, errors.New("order not found")
	}
	_, err = tx.Updates(o)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (r *orderRepository) GetOrderByStatus(ctx context.Context, status string) ([]*model.Order, error) {
	tx := newOrder(r.DB(ctx))
	return tx.Where(tx.Status.Eq(status)).Limit(50).Find()
}
