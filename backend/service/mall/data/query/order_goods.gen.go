// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"fastduck/treasure-doc/service/mall/data/model"
)

func newOrderGood(db *gorm.DB, opts ...gen.DOOption) orderGood {
	_orderGood := orderGood{}

	_orderGood.orderGoodDo.UseDB(db, opts...)
	_orderGood.orderGoodDo.UseModel(&model.OrderGood{})

	tableName := _orderGood.orderGoodDo.TableName()
	_orderGood.ALL = field.NewAsterisk(tableName)
	_orderGood.ID = field.NewInt32(tableName, "id")
	_orderGood.GoodID = field.NewInt32(tableName, "good_id")
	_orderGood.CreatedAt = field.NewTime(tableName, "created_at")
	_orderGood.UpdatedAt = field.NewTime(tableName, "updated_at")
	_orderGood.DeletedAt = field.NewField(tableName, "deleted_at")

	_orderGood.fillFieldMap()

	return _orderGood
}

type orderGood struct {
	orderGoodDo

	ALL       field.Asterisk
	ID        field.Int32
	GoodID    field.Int32
	CreatedAt field.Time
	UpdatedAt field.Time
	DeletedAt field.Field

	fieldMap map[string]field.Expr
}

func (o orderGood) Table(newTableName string) *orderGood {
	o.orderGoodDo.UseTable(newTableName)
	return o.updateTableName(newTableName)
}

func (o orderGood) As(alias string) *orderGood {
	o.orderGoodDo.DO = *(o.orderGoodDo.As(alias).(*gen.DO))
	return o.updateTableName(alias)
}

func (o *orderGood) updateTableName(table string) *orderGood {
	o.ALL = field.NewAsterisk(table)
	o.ID = field.NewInt32(table, "id")
	o.GoodID = field.NewInt32(table, "good_id")
	o.CreatedAt = field.NewTime(table, "created_at")
	o.UpdatedAt = field.NewTime(table, "updated_at")
	o.DeletedAt = field.NewField(table, "deleted_at")

	o.fillFieldMap()

	return o
}

func (o *orderGood) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := o.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (o *orderGood) fillFieldMap() {
	o.fieldMap = make(map[string]field.Expr, 5)
	o.fieldMap["id"] = o.ID
	o.fieldMap["good_id"] = o.GoodID
	o.fieldMap["created_at"] = o.CreatedAt
	o.fieldMap["updated_at"] = o.UpdatedAt
	o.fieldMap["deleted_at"] = o.DeletedAt
}

func (o orderGood) clone(db *gorm.DB) orderGood {
	o.orderGoodDo.ReplaceConnPool(db.Statement.ConnPool)
	return o
}

func (o orderGood) replaceDB(db *gorm.DB) orderGood {
	o.orderGoodDo.ReplaceDB(db)
	return o
}

type orderGoodDo struct{ gen.DO }

type IOrderGoodDo interface {
	gen.SubQuery
	Debug() IOrderGoodDo
	WithContext(ctx context.Context) IOrderGoodDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IOrderGoodDo
	WriteDB() IOrderGoodDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IOrderGoodDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IOrderGoodDo
	Not(conds ...gen.Condition) IOrderGoodDo
	Or(conds ...gen.Condition) IOrderGoodDo
	Select(conds ...field.Expr) IOrderGoodDo
	Where(conds ...gen.Condition) IOrderGoodDo
	Order(conds ...field.Expr) IOrderGoodDo
	Distinct(cols ...field.Expr) IOrderGoodDo
	Omit(cols ...field.Expr) IOrderGoodDo
	Join(table schema.Tabler, on ...field.Expr) IOrderGoodDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IOrderGoodDo
	RightJoin(table schema.Tabler, on ...field.Expr) IOrderGoodDo
	Group(cols ...field.Expr) IOrderGoodDo
	Having(conds ...gen.Condition) IOrderGoodDo
	Limit(limit int) IOrderGoodDo
	Offset(offset int) IOrderGoodDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IOrderGoodDo
	Unscoped() IOrderGoodDo
	Create(values ...*model.OrderGood) error
	CreateInBatches(values []*model.OrderGood, batchSize int) error
	Save(values ...*model.OrderGood) error
	First() (*model.OrderGood, error)
	Take() (*model.OrderGood, error)
	Last() (*model.OrderGood, error)
	Find() ([]*model.OrderGood, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.OrderGood, err error)
	FindInBatches(result *[]*model.OrderGood, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.OrderGood) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IOrderGoodDo
	Assign(attrs ...field.AssignExpr) IOrderGoodDo
	Joins(fields ...field.RelationField) IOrderGoodDo
	Preload(fields ...field.RelationField) IOrderGoodDo
	FirstOrInit() (*model.OrderGood, error)
	FirstOrCreate() (*model.OrderGood, error)
	FindByPage(offset int, limit int) (result []*model.OrderGood, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IOrderGoodDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (o orderGoodDo) Debug() IOrderGoodDo {
	return o.withDO(o.DO.Debug())
}

func (o orderGoodDo) WithContext(ctx context.Context) IOrderGoodDo {
	return o.withDO(o.DO.WithContext(ctx))
}

func (o orderGoodDo) ReadDB() IOrderGoodDo {
	return o.Clauses(dbresolver.Read)
}

func (o orderGoodDo) WriteDB() IOrderGoodDo {
	return o.Clauses(dbresolver.Write)
}

func (o orderGoodDo) Session(config *gorm.Session) IOrderGoodDo {
	return o.withDO(o.DO.Session(config))
}

func (o orderGoodDo) Clauses(conds ...clause.Expression) IOrderGoodDo {
	return o.withDO(o.DO.Clauses(conds...))
}

func (o orderGoodDo) Returning(value interface{}, columns ...string) IOrderGoodDo {
	return o.withDO(o.DO.Returning(value, columns...))
}

func (o orderGoodDo) Not(conds ...gen.Condition) IOrderGoodDo {
	return o.withDO(o.DO.Not(conds...))
}

func (o orderGoodDo) Or(conds ...gen.Condition) IOrderGoodDo {
	return o.withDO(o.DO.Or(conds...))
}

func (o orderGoodDo) Select(conds ...field.Expr) IOrderGoodDo {
	return o.withDO(o.DO.Select(conds...))
}

func (o orderGoodDo) Where(conds ...gen.Condition) IOrderGoodDo {
	return o.withDO(o.DO.Where(conds...))
}

func (o orderGoodDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) IOrderGoodDo {
	return o.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (o orderGoodDo) Order(conds ...field.Expr) IOrderGoodDo {
	return o.withDO(o.DO.Order(conds...))
}

func (o orderGoodDo) Distinct(cols ...field.Expr) IOrderGoodDo {
	return o.withDO(o.DO.Distinct(cols...))
}

func (o orderGoodDo) Omit(cols ...field.Expr) IOrderGoodDo {
	return o.withDO(o.DO.Omit(cols...))
}

func (o orderGoodDo) Join(table schema.Tabler, on ...field.Expr) IOrderGoodDo {
	return o.withDO(o.DO.Join(table, on...))
}

func (o orderGoodDo) LeftJoin(table schema.Tabler, on ...field.Expr) IOrderGoodDo {
	return o.withDO(o.DO.LeftJoin(table, on...))
}

func (o orderGoodDo) RightJoin(table schema.Tabler, on ...field.Expr) IOrderGoodDo {
	return o.withDO(o.DO.RightJoin(table, on...))
}

func (o orderGoodDo) Group(cols ...field.Expr) IOrderGoodDo {
	return o.withDO(o.DO.Group(cols...))
}

func (o orderGoodDo) Having(conds ...gen.Condition) IOrderGoodDo {
	return o.withDO(o.DO.Having(conds...))
}

func (o orderGoodDo) Limit(limit int) IOrderGoodDo {
	return o.withDO(o.DO.Limit(limit))
}

func (o orderGoodDo) Offset(offset int) IOrderGoodDo {
	return o.withDO(o.DO.Offset(offset))
}

func (o orderGoodDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IOrderGoodDo {
	return o.withDO(o.DO.Scopes(funcs...))
}

func (o orderGoodDo) Unscoped() IOrderGoodDo {
	return o.withDO(o.DO.Unscoped())
}

func (o orderGoodDo) Create(values ...*model.OrderGood) error {
	if len(values) == 0 {
		return nil
	}
	return o.DO.Create(values)
}

func (o orderGoodDo) CreateInBatches(values []*model.OrderGood, batchSize int) error {
	return o.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (o orderGoodDo) Save(values ...*model.OrderGood) error {
	if len(values) == 0 {
		return nil
	}
	return o.DO.Save(values)
}

func (o orderGoodDo) First() (*model.OrderGood, error) {
	if result, err := o.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.OrderGood), nil
	}
}

func (o orderGoodDo) Take() (*model.OrderGood, error) {
	if result, err := o.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.OrderGood), nil
	}
}

func (o orderGoodDo) Last() (*model.OrderGood, error) {
	if result, err := o.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.OrderGood), nil
	}
}

func (o orderGoodDo) Find() ([]*model.OrderGood, error) {
	result, err := o.DO.Find()
	return result.([]*model.OrderGood), err
}

func (o orderGoodDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.OrderGood, err error) {
	buf := make([]*model.OrderGood, 0, batchSize)
	err = o.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (o orderGoodDo) FindInBatches(result *[]*model.OrderGood, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return o.DO.FindInBatches(result, batchSize, fc)
}

func (o orderGoodDo) Attrs(attrs ...field.AssignExpr) IOrderGoodDo {
	return o.withDO(o.DO.Attrs(attrs...))
}

func (o orderGoodDo) Assign(attrs ...field.AssignExpr) IOrderGoodDo {
	return o.withDO(o.DO.Assign(attrs...))
}

func (o orderGoodDo) Joins(fields ...field.RelationField) IOrderGoodDo {
	for _, _f := range fields {
		o = *o.withDO(o.DO.Joins(_f))
	}
	return &o
}

func (o orderGoodDo) Preload(fields ...field.RelationField) IOrderGoodDo {
	for _, _f := range fields {
		o = *o.withDO(o.DO.Preload(_f))
	}
	return &o
}

func (o orderGoodDo) FirstOrInit() (*model.OrderGood, error) {
	if result, err := o.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.OrderGood), nil
	}
}

func (o orderGoodDo) FirstOrCreate() (*model.OrderGood, error) {
	if result, err := o.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.OrderGood), nil
	}
}

func (o orderGoodDo) FindByPage(offset int, limit int) (result []*model.OrderGood, count int64, err error) {
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

func (o orderGoodDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = o.Count()
	if err != nil {
		return
	}

	err = o.Offset(offset).Limit(limit).Scan(result)
	return
}

func (o orderGoodDo) Scan(result interface{}) (err error) {
	return o.DO.Scan(result)
}

func (o orderGoodDo) Delete(models ...*model.OrderGood) (result gen.ResultInfo, err error) {
	return o.DO.Delete(models)
}

func (o *orderGoodDo) withDO(do gen.Dao) *orderGoodDo {
	o.DO = *do.(*gen.DO)
	return o
}
