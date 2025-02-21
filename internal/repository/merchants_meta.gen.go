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

func newMerchantsMetum(db *gorm.DB, opts ...gen.DOOption) merchantsMetum {
	_merchantsMetum := merchantsMetum{}

	_merchantsMetum.merchantsMetumDo.UseDB(db, opts...)
	_merchantsMetum.merchantsMetumDo.UseModel(&model.MerchantsMetum{})

	tableName := _merchantsMetum.merchantsMetumDo.TableName()
	_merchantsMetum.ALL = field.NewAsterisk(tableName)
	_merchantsMetum.ID = field.NewInt64(tableName, "id")
	_merchantsMetum.MID = field.NewInt64(tableName, "m_id")
	_merchantsMetum.Ac = field.NewString(tableName, "ac")
	_merchantsMetum.CreatedAt = field.NewTime(tableName, "created_at")
	_merchantsMetum.UpdatedAt = field.NewTime(tableName, "updated_at")
	_merchantsMetum.DeletedAt = field.NewField(tableName, "deleted_at")

	_merchantsMetum.fillFieldMap()

	return _merchantsMetum
}

type merchantsMetum struct {
	merchantsMetumDo

	ALL       field.Asterisk
	ID        field.Int64
	MID       field.Int64
	Ac        field.String // 钱包地址
	CreatedAt field.Time
	UpdatedAt field.Time
	DeletedAt field.Field

	fieldMap map[string]field.Expr
}

func (m merchantsMetum) Table(newTableName string) *merchantsMetum {
	m.merchantsMetumDo.UseTable(newTableName)
	return m.updateTableName(newTableName)
}

func (m merchantsMetum) As(alias string) *merchantsMetum {
	m.merchantsMetumDo.DO = *(m.merchantsMetumDo.As(alias).(*gen.DO))
	return m.updateTableName(alias)
}

func (m *merchantsMetum) updateTableName(table string) *merchantsMetum {
	m.ALL = field.NewAsterisk(table)
	m.ID = field.NewInt64(table, "id")
	m.MID = field.NewInt64(table, "m_id")
	m.Ac = field.NewString(table, "ac")
	m.CreatedAt = field.NewTime(table, "created_at")
	m.UpdatedAt = field.NewTime(table, "updated_at")
	m.DeletedAt = field.NewField(table, "deleted_at")

	m.fillFieldMap()

	return m
}

func (m *merchantsMetum) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := m.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (m *merchantsMetum) fillFieldMap() {
	m.fieldMap = make(map[string]field.Expr, 6)
	m.fieldMap["id"] = m.ID
	m.fieldMap["m_id"] = m.MID
	m.fieldMap["ac"] = m.Ac
	m.fieldMap["created_at"] = m.CreatedAt
	m.fieldMap["updated_at"] = m.UpdatedAt
	m.fieldMap["deleted_at"] = m.DeletedAt
}

func (m merchantsMetum) clone(db *gorm.DB) merchantsMetum {
	m.merchantsMetumDo.ReplaceConnPool(db.Statement.ConnPool)
	return m
}

func (m merchantsMetum) replaceDB(db *gorm.DB) merchantsMetum {
	m.merchantsMetumDo.ReplaceDB(db)
	return m
}

type merchantsMetumDo struct{ gen.DO }

type IMerchantsMetumDo interface {
	gen.SubQuery
	Debug() IMerchantsMetumDo
	WithContext(ctx context.Context) IMerchantsMetumDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IMerchantsMetumDo
	WriteDB() IMerchantsMetumDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IMerchantsMetumDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IMerchantsMetumDo
	Not(conds ...gen.Condition) IMerchantsMetumDo
	Or(conds ...gen.Condition) IMerchantsMetumDo
	Select(conds ...field.Expr) IMerchantsMetumDo
	Where(conds ...gen.Condition) IMerchantsMetumDo
	Order(conds ...field.Expr) IMerchantsMetumDo
	Distinct(cols ...field.Expr) IMerchantsMetumDo
	Omit(cols ...field.Expr) IMerchantsMetumDo
	Join(table schema.Tabler, on ...field.Expr) IMerchantsMetumDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IMerchantsMetumDo
	RightJoin(table schema.Tabler, on ...field.Expr) IMerchantsMetumDo
	Group(cols ...field.Expr) IMerchantsMetumDo
	Having(conds ...gen.Condition) IMerchantsMetumDo
	Limit(limit int) IMerchantsMetumDo
	Offset(offset int) IMerchantsMetumDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IMerchantsMetumDo
	Unscoped() IMerchantsMetumDo
	Create(values ...*model.MerchantsMetum) error
	CreateInBatches(values []*model.MerchantsMetum, batchSize int) error
	Save(values ...*model.MerchantsMetum) error
	First() (*model.MerchantsMetum, error)
	Take() (*model.MerchantsMetum, error)
	Last() (*model.MerchantsMetum, error)
	Find() ([]*model.MerchantsMetum, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.MerchantsMetum, err error)
	FindInBatches(result *[]*model.MerchantsMetum, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.MerchantsMetum) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IMerchantsMetumDo
	Assign(attrs ...field.AssignExpr) IMerchantsMetumDo
	Joins(fields ...field.RelationField) IMerchantsMetumDo
	Preload(fields ...field.RelationField) IMerchantsMetumDo
	FirstOrInit() (*model.MerchantsMetum, error)
	FirstOrCreate() (*model.MerchantsMetum, error)
	FindByPage(offset int, limit int) (result []*model.MerchantsMetum, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IMerchantsMetumDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (m merchantsMetumDo) Debug() IMerchantsMetumDo {
	return m.withDO(m.DO.Debug())
}

func (m merchantsMetumDo) WithContext(ctx context.Context) IMerchantsMetumDo {
	return m.withDO(m.DO.WithContext(ctx))
}

func (m merchantsMetumDo) ReadDB() IMerchantsMetumDo {
	return m.Clauses(dbresolver.Read)
}

func (m merchantsMetumDo) WriteDB() IMerchantsMetumDo {
	return m.Clauses(dbresolver.Write)
}

func (m merchantsMetumDo) Session(config *gorm.Session) IMerchantsMetumDo {
	return m.withDO(m.DO.Session(config))
}

func (m merchantsMetumDo) Clauses(conds ...clause.Expression) IMerchantsMetumDo {
	return m.withDO(m.DO.Clauses(conds...))
}

func (m merchantsMetumDo) Returning(value interface{}, columns ...string) IMerchantsMetumDo {
	return m.withDO(m.DO.Returning(value, columns...))
}

func (m merchantsMetumDo) Not(conds ...gen.Condition) IMerchantsMetumDo {
	return m.withDO(m.DO.Not(conds...))
}

func (m merchantsMetumDo) Or(conds ...gen.Condition) IMerchantsMetumDo {
	return m.withDO(m.DO.Or(conds...))
}

func (m merchantsMetumDo) Select(conds ...field.Expr) IMerchantsMetumDo {
	return m.withDO(m.DO.Select(conds...))
}

func (m merchantsMetumDo) Where(conds ...gen.Condition) IMerchantsMetumDo {
	return m.withDO(m.DO.Where(conds...))
}

func (m merchantsMetumDo) Order(conds ...field.Expr) IMerchantsMetumDo {
	return m.withDO(m.DO.Order(conds...))
}

func (m merchantsMetumDo) Distinct(cols ...field.Expr) IMerchantsMetumDo {
	return m.withDO(m.DO.Distinct(cols...))
}

func (m merchantsMetumDo) Omit(cols ...field.Expr) IMerchantsMetumDo {
	return m.withDO(m.DO.Omit(cols...))
}

func (m merchantsMetumDo) Join(table schema.Tabler, on ...field.Expr) IMerchantsMetumDo {
	return m.withDO(m.DO.Join(table, on...))
}

func (m merchantsMetumDo) LeftJoin(table schema.Tabler, on ...field.Expr) IMerchantsMetumDo {
	return m.withDO(m.DO.LeftJoin(table, on...))
}

func (m merchantsMetumDo) RightJoin(table schema.Tabler, on ...field.Expr) IMerchantsMetumDo {
	return m.withDO(m.DO.RightJoin(table, on...))
}

func (m merchantsMetumDo) Group(cols ...field.Expr) IMerchantsMetumDo {
	return m.withDO(m.DO.Group(cols...))
}

func (m merchantsMetumDo) Having(conds ...gen.Condition) IMerchantsMetumDo {
	return m.withDO(m.DO.Having(conds...))
}

func (m merchantsMetumDo) Limit(limit int) IMerchantsMetumDo {
	return m.withDO(m.DO.Limit(limit))
}

func (m merchantsMetumDo) Offset(offset int) IMerchantsMetumDo {
	return m.withDO(m.DO.Offset(offset))
}

func (m merchantsMetumDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IMerchantsMetumDo {
	return m.withDO(m.DO.Scopes(funcs...))
}

func (m merchantsMetumDo) Unscoped() IMerchantsMetumDo {
	return m.withDO(m.DO.Unscoped())
}

func (m merchantsMetumDo) Create(values ...*model.MerchantsMetum) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Create(values)
}

func (m merchantsMetumDo) CreateInBatches(values []*model.MerchantsMetum, batchSize int) error {
	return m.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (m merchantsMetumDo) Save(values ...*model.MerchantsMetum) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Save(values)
}

func (m merchantsMetumDo) First() (*model.MerchantsMetum, error) {
	if result, err := m.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.MerchantsMetum), nil
	}
}

func (m merchantsMetumDo) Take() (*model.MerchantsMetum, error) {
	if result, err := m.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.MerchantsMetum), nil
	}
}

func (m merchantsMetumDo) Last() (*model.MerchantsMetum, error) {
	if result, err := m.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.MerchantsMetum), nil
	}
}

func (m merchantsMetumDo) Find() ([]*model.MerchantsMetum, error) {
	result, err := m.DO.Find()
	return result.([]*model.MerchantsMetum), err
}

func (m merchantsMetumDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.MerchantsMetum, err error) {
	buf := make([]*model.MerchantsMetum, 0, batchSize)
	err = m.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (m merchantsMetumDo) FindInBatches(result *[]*model.MerchantsMetum, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return m.DO.FindInBatches(result, batchSize, fc)
}

func (m merchantsMetumDo) Attrs(attrs ...field.AssignExpr) IMerchantsMetumDo {
	return m.withDO(m.DO.Attrs(attrs...))
}

func (m merchantsMetumDo) Assign(attrs ...field.AssignExpr) IMerchantsMetumDo {
	return m.withDO(m.DO.Assign(attrs...))
}

func (m merchantsMetumDo) Joins(fields ...field.RelationField) IMerchantsMetumDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Joins(_f))
	}
	return &m
}

func (m merchantsMetumDo) Preload(fields ...field.RelationField) IMerchantsMetumDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Preload(_f))
	}
	return &m
}

func (m merchantsMetumDo) FirstOrInit() (*model.MerchantsMetum, error) {
	if result, err := m.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.MerchantsMetum), nil
	}
}

func (m merchantsMetumDo) FirstOrCreate() (*model.MerchantsMetum, error) {
	if result, err := m.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.MerchantsMetum), nil
	}
}

func (m merchantsMetumDo) FindByPage(offset int, limit int) (result []*model.MerchantsMetum, count int64, err error) {
	result, err = m.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = m.Offset(-1).Limit(-1).Count()
	return
}

func (m merchantsMetumDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = m.Count()
	if err != nil {
		return
	}

	err = m.Offset(offset).Limit(limit).Scan(result)
	return
}

func (m merchantsMetumDo) Scan(result interface{}) (err error) {
	return m.DO.Scan(result)
}

func (m merchantsMetumDo) Delete(models ...*model.MerchantsMetum) (result gen.ResultInfo, err error) {
	return m.DO.Delete(models)
}

func (m *merchantsMetumDo) withDO(do gen.Dao) *merchantsMetumDo {
	m.DO = *do.(*gen.DO)
	return m
}
