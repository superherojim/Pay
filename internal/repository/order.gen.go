// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package repository

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"cheemshappy_pay/internal/model"
)

func newOrder(db *gorm.DB, opts ...gen.DOOption) order {
	_order := order{}

	_order.orderDo.UseDB(db, opts...)
	_order.orderDo.UseModel(&model.Order{})

	tableName := _order.orderDo.TableName()
	_order.ALL = field.NewAsterisk(tableName)
	_order.ID = field.NewInt64(tableName, "id")
	_order.MID = field.NewInt64(tableName, "m_id")
	_order.No = field.NewString(tableName, "no")
	_order.CNo = field.NewString(tableName, "c_no")
	_order.Amount = field.NewString(tableName, "amount")
	_order.Coin = field.NewString(tableName, "coin")
	_order.Chain = field.NewString(tableName, "chain")
	_order.NotifyURL = field.NewString(tableName, "notify_url")
	_order.ReturnURL = field.NewString(tableName, "return_url")
	_order.PayURL = field.NewString(tableName, "pay_url")
	_order.TxHash = field.NewString(tableName, "tx_hash")
	_order.Ac = field.NewString(tableName, "ac")
	_order.APIKey = field.NewString(tableName, "api_key")
	_order.TimeOut = field.NewInt32(tableName, "time_out")
	_order.Status = field.NewString(tableName, "status")
	_order.NotifyStatus = field.NewString(tableName, "notify_status")
	_order.Remark = field.NewString(tableName, "remark")
	_order.CreatedAt = field.NewTime(tableName, "created_at")
	_order.UpdatedAt = field.NewTime(tableName, "updated_at")
	_order.DeletedAt = field.NewField(tableName, "deleted_at")

	_order.fillFieldMap()

	return _order
}

type order struct {
	orderDo

	ALL          field.Asterisk
	ID           field.Int64
	MID          field.Int64
	No           field.String // 订单no
	CNo          field.String // 用户no
	Amount       field.String // 价格
	Coin         field.String // 数字币类型
	Chain        field.String // 链
	NotifyURL    field.String // 回调url
	ReturnURL    field.String // 前端重定向url
	PayURL       field.String // 支付地址
	TxHash       field.String // 交易hash
	Ac           field.String // 钱包地址
	APIKey       field.String // 商户apikey
	TimeOut      field.Int32  // 超时时间
	Status       field.String // 订单状态
	NotifyStatus field.String // 订单回调结果
	Remark       field.String // 备注
	CreatedAt    field.Time
	UpdatedAt    field.Time
	DeletedAt    field.Field

	fieldMap map[string]field.Expr
}

func (o order) Table(newTableName string) *order {
	o.orderDo.UseTable(newTableName)
	return o.updateTableName(newTableName)
}

func (o order) As(alias string) *order {
	o.orderDo.DO = *(o.orderDo.As(alias).(*gen.DO))
	return o.updateTableName(alias)
}

func (o *order) updateTableName(table string) *order {
	o.ALL = field.NewAsterisk(table)
	o.ID = field.NewInt64(table, "id")
	o.MID = field.NewInt64(table, "m_id")
	o.No = field.NewString(table, "no")
	o.CNo = field.NewString(table, "c_no")
	o.Amount = field.NewString(table, "amount")
	o.Coin = field.NewString(table, "coin")
	o.Chain = field.NewString(table, "chain")
	o.NotifyURL = field.NewString(table, "notify_url")
	o.ReturnURL = field.NewString(table, "return_url")
	o.PayURL = field.NewString(table, "pay_url")
	o.TxHash = field.NewString(table, "tx_hash")
	o.Ac = field.NewString(table, "ac")
	o.APIKey = field.NewString(table, "api_key")
	o.TimeOut = field.NewInt32(table, "time_out")
	o.Status = field.NewString(table, "status")
	o.NotifyStatus = field.NewString(table, "notify_status")
	o.Remark = field.NewString(table, "remark")
	o.CreatedAt = field.NewTime(table, "created_at")
	o.UpdatedAt = field.NewTime(table, "updated_at")
	o.DeletedAt = field.NewField(table, "deleted_at")

	o.fillFieldMap()

	return o
}

func (o *order) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := o.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (o *order) fillFieldMap() {
	o.fieldMap = make(map[string]field.Expr, 20)
	o.fieldMap["id"] = o.ID
	o.fieldMap["m_id"] = o.MID
	o.fieldMap["no"] = o.No
	o.fieldMap["c_no"] = o.CNo
	o.fieldMap["amount"] = o.Amount
	o.fieldMap["coin"] = o.Coin
	o.fieldMap["chain"] = o.Chain
	o.fieldMap["notify_url"] = o.NotifyURL
	o.fieldMap["return_url"] = o.ReturnURL
	o.fieldMap["pay_url"] = o.PayURL
	o.fieldMap["tx_hash"] = o.TxHash
	o.fieldMap["ac"] = o.Ac
	o.fieldMap["api_key"] = o.APIKey
	o.fieldMap["time_out"] = o.TimeOut
	o.fieldMap["status"] = o.Status
	o.fieldMap["notify_status"] = o.NotifyStatus
	o.fieldMap["remark"] = o.Remark
	o.fieldMap["created_at"] = o.CreatedAt
	o.fieldMap["updated_at"] = o.UpdatedAt
	o.fieldMap["deleted_at"] = o.DeletedAt
}

func (o order) clone(db *gorm.DB) order {
	o.orderDo.ReplaceConnPool(db.Statement.ConnPool)
	return o
}

func (o order) replaceDB(db *gorm.DB) order {
	o.orderDo.ReplaceDB(db)
	return o
}

type orderDo struct{ gen.DO }

type IOrderDo interface {
	gen.SubQuery
	Debug() IOrderDo
	WithContext(ctx context.Context) IOrderDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IOrderDo
	WriteDB() IOrderDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IOrderDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IOrderDo
	Not(conds ...gen.Condition) IOrderDo
	Or(conds ...gen.Condition) IOrderDo
	Select(conds ...field.Expr) IOrderDo
	Where(conds ...gen.Condition) IOrderDo
	Order(conds ...field.Expr) IOrderDo
	Distinct(cols ...field.Expr) IOrderDo
	Omit(cols ...field.Expr) IOrderDo
	Join(table schema.Tabler, on ...field.Expr) IOrderDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IOrderDo
	RightJoin(table schema.Tabler, on ...field.Expr) IOrderDo
	Group(cols ...field.Expr) IOrderDo
	Having(conds ...gen.Condition) IOrderDo
	Limit(limit int) IOrderDo
	Offset(offset int) IOrderDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IOrderDo
	Unscoped() IOrderDo
	Create(values ...*model.Order) error
	CreateInBatches(values []*model.Order, batchSize int) error
	Save(values ...*model.Order) error
	First() (*model.Order, error)
	Take() (*model.Order, error)
	Last() (*model.Order, error)
	Find() ([]*model.Order, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Order, err error)
	FindInBatches(result *[]*model.Order, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.Order) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IOrderDo
	Assign(attrs ...field.AssignExpr) IOrderDo
	Joins(fields ...field.RelationField) IOrderDo
	Preload(fields ...field.RelationField) IOrderDo
	FirstOrInit() (*model.Order, error)
	FirstOrCreate() (*model.Order, error)
	FindByPage(offset int, limit int) (result []*model.Order, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IOrderDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (o orderDo) Debug() IOrderDo {
	return o.withDO(o.DO.Debug())
}

func (o orderDo) WithContext(ctx context.Context) IOrderDo {
	return o.withDO(o.DO.WithContext(ctx))
}

func (o orderDo) ReadDB() IOrderDo {
	return o.Clauses(dbresolver.Read)
}

func (o orderDo) WriteDB() IOrderDo {
	return o.Clauses(dbresolver.Write)
}

func (o orderDo) Session(config *gorm.Session) IOrderDo {
	return o.withDO(o.DO.Session(config))
}

func (o orderDo) Clauses(conds ...clause.Expression) IOrderDo {
	return o.withDO(o.DO.Clauses(conds...))
}

func (o orderDo) Returning(value interface{}, columns ...string) IOrderDo {
	return o.withDO(o.DO.Returning(value, columns...))
}

func (o orderDo) Not(conds ...gen.Condition) IOrderDo {
	return o.withDO(o.DO.Not(conds...))
}

func (o orderDo) Or(conds ...gen.Condition) IOrderDo {
	return o.withDO(o.DO.Or(conds...))
}

func (o orderDo) Select(conds ...field.Expr) IOrderDo {
	return o.withDO(o.DO.Select(conds...))
}

func (o orderDo) Where(conds ...gen.Condition) IOrderDo {
	return o.withDO(o.DO.Where(conds...))
}

func (o orderDo) Order(conds ...field.Expr) IOrderDo {
	return o.withDO(o.DO.Order(conds...))
}

func (o orderDo) Distinct(cols ...field.Expr) IOrderDo {
	return o.withDO(o.DO.Distinct(cols...))
}

func (o orderDo) Omit(cols ...field.Expr) IOrderDo {
	return o.withDO(o.DO.Omit(cols...))
}

func (o orderDo) Join(table schema.Tabler, on ...field.Expr) IOrderDo {
	return o.withDO(o.DO.Join(table, on...))
}

func (o orderDo) LeftJoin(table schema.Tabler, on ...field.Expr) IOrderDo {
	return o.withDO(o.DO.LeftJoin(table, on...))
}

func (o orderDo) RightJoin(table schema.Tabler, on ...field.Expr) IOrderDo {
	return o.withDO(o.DO.RightJoin(table, on...))
}

func (o orderDo) Group(cols ...field.Expr) IOrderDo {
	return o.withDO(o.DO.Group(cols...))
}

func (o orderDo) Having(conds ...gen.Condition) IOrderDo {
	return o.withDO(o.DO.Having(conds...))
}

func (o orderDo) Limit(limit int) IOrderDo {
	return o.withDO(o.DO.Limit(limit))
}

func (o orderDo) Offset(offset int) IOrderDo {
	return o.withDO(o.DO.Offset(offset))
}

func (o orderDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IOrderDo {
	return o.withDO(o.DO.Scopes(funcs...))
}

func (o orderDo) Unscoped() IOrderDo {
	return o.withDO(o.DO.Unscoped())
}

func (o orderDo) Create(values ...*model.Order) error {
	if len(values) == 0 {
		return nil
	}
	return o.DO.Create(values)
}

func (o orderDo) CreateInBatches(values []*model.Order, batchSize int) error {
	return o.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (o orderDo) Save(values ...*model.Order) error {
	if len(values) == 0 {
		return nil
	}
	return o.DO.Save(values)
}

func (o orderDo) First() (*model.Order, error) {
	if result, err := o.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Order), nil
	}
}

func (o orderDo) Take() (*model.Order, error) {
	if result, err := o.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Order), nil
	}
}

func (o orderDo) Last() (*model.Order, error) {
	if result, err := o.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Order), nil
	}
}

func (o orderDo) Find() ([]*model.Order, error) {
	result, err := o.DO.Find()
	return result.([]*model.Order), err
}

func (o orderDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Order, err error) {
	buf := make([]*model.Order, 0, batchSize)
	err = o.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (o orderDo) FindInBatches(result *[]*model.Order, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return o.DO.FindInBatches(result, batchSize, fc)
}

func (o orderDo) Attrs(attrs ...field.AssignExpr) IOrderDo {
	return o.withDO(o.DO.Attrs(attrs...))
}

func (o orderDo) Assign(attrs ...field.AssignExpr) IOrderDo {
	return o.withDO(o.DO.Assign(attrs...))
}

func (o orderDo) Joins(fields ...field.RelationField) IOrderDo {
	for _, _f := range fields {
		o = *o.withDO(o.DO.Joins(_f))
	}
	return &o
}

func (o orderDo) Preload(fields ...field.RelationField) IOrderDo {
	for _, _f := range fields {
		o = *o.withDO(o.DO.Preload(_f))
	}
	return &o
}

func (o orderDo) FirstOrInit() (*model.Order, error) {
	if result, err := o.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Order), nil
	}
}

func (o orderDo) FirstOrCreate() (*model.Order, error) {
	if result, err := o.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Order), nil
	}
}

func (o orderDo) FindByPage(offset int, limit int) (result []*model.Order, count int64, err error) {
	result, err = o.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = o.Offset(-1).Limit(-1).Count()
	return
}

func (o orderDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = o.Count()
	if err != nil {
		return
	}

	err = o.Offset(offset).Limit(limit).Scan(result)
	return
}

func (o orderDo) Scan(result interface{}) (err error) {
	return o.DO.Scan(result)
}

func (o orderDo) Delete(models ...*model.Order) (result gen.ResultInfo, err error) {
	return o.DO.Delete(models)
}

func (o *orderDo) withDO(do gen.Dao) *orderDo {
	o.DO = *do.(*gen.DO)
	return o
}
