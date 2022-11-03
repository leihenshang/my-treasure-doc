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

func newOrderDetail(db *gorm.DB, opts ...gen.DOOption) orderDetail {
	_orderDetail := orderDetail{}

	_orderDetail.orderDetailDo.UseDB(db, opts...)
	_orderDetail.orderDetailDo.UseModel(&model.OrderDetail{})

	tableName := _orderDetail.orderDetailDo.TableName()
	_orderDetail.ALL = field.NewAsterisk(tableName)
	_orderDetail.ID = field.NewInt32(tableName, "id")
	_orderDetail.OrderID = field.NewInt32(tableName, "order_id")
	_orderDetail.GoodID = field.NewInt32(tableName, "good_id")
	_orderDetail.SkuID = field.NewInt32(tableName, "sku_id")
	_orderDetail.Price = field.NewFloat64(tableName, "price")
	_orderDetail.Quantity = field.NewInt32(tableName, "quantity")
	_orderDetail.CreatedAt = field.NewTime(tableName, "created_at")
	_orderDetail.UpdatedAt = field.NewTime(tableName, "updated_at")
	_orderDetail.DeletedAt = field.NewField(tableName, "deleted_at")

	_orderDetail.fillFieldMap()

	return _orderDetail
}

type orderDetail struct {
	orderDetailDo

	ALL       field.Asterisk
	ID        field.Int32
	OrderID   field.Int32 // 订单id
	GoodID    field.Int32
	SkuID     field.Int32   // sku id
	Price     field.Float64 // 单价
	Quantity  field.Int32   // 数量
	CreatedAt field.Time
	UpdatedAt field.Time
	DeletedAt field.Field

	fieldMap map[string]field.Expr
}

func (o orderDetail) Table(newTableName string) *orderDetail {
	o.orderDetailDo.UseTable(newTableName)
	return o.updateTableName(newTableName)
}

func (o orderDetail) As(alias string) *orderDetail {
	o.orderDetailDo.DO = *(o.orderDetailDo.As(alias).(*gen.DO))
	return o.updateTableName(alias)
}

func (o *orderDetail) updateTableName(table string) *orderDetail {
	o.ALL = field.NewAsterisk(table)
	o.ID = field.NewInt32(table, "id")
	o.OrderID = field.NewInt32(table, "order_id")
	o.GoodID = field.NewInt32(table, "good_id")
	o.SkuID = field.NewInt32(table, "sku_id")
	o.Price = field.NewFloat64(table, "price")
	o.Quantity = field.NewInt32(table, "quantity")
	o.CreatedAt = field.NewTime(table, "created_at")
	o.UpdatedAt = field.NewTime(table, "updated_at")
	o.DeletedAt = field.NewField(table, "deleted_at")

	o.fillFieldMap()

	return o
}

func (o *orderDetail) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := o.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (o *orderDetail) fillFieldMap() {
	o.fieldMap = make(map[string]field.Expr, 9)
	o.fieldMap["id"] = o.ID
	o.fieldMap["order_id"] = o.OrderID
	o.fieldMap["good_id"] = o.GoodID
	o.fieldMap["sku_id"] = o.SkuID
	o.fieldMap["price"] = o.Price
	o.fieldMap["quantity"] = o.Quantity
	o.fieldMap["created_at"] = o.CreatedAt
	o.fieldMap["updated_at"] = o.UpdatedAt
	o.fieldMap["deleted_at"] = o.DeletedAt
}

func (o orderDetail) clone(db *gorm.DB) orderDetail {
	o.orderDetailDo.ReplaceConnPool(db.Statement.ConnPool)
	return o
}

func (o orderDetail) replaceDB(db *gorm.DB) orderDetail {
	o.orderDetailDo.ReplaceDB(db)
	return o
}

type orderDetailDo struct{ gen.DO }

type IOrderDetailDo interface {
	gen.SubQuery
	Debug() IOrderDetailDo
	WithContext(ctx context.Context) IOrderDetailDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IOrderDetailDo
	WriteDB() IOrderDetailDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IOrderDetailDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IOrderDetailDo
	Not(conds ...gen.Condition) IOrderDetailDo
	Or(conds ...gen.Condition) IOrderDetailDo
	Select(conds ...field.Expr) IOrderDetailDo
	Where(conds ...gen.Condition) IOrderDetailDo
	Order(conds ...field.Expr) IOrderDetailDo
	Distinct(cols ...field.Expr) IOrderDetailDo
	Omit(cols ...field.Expr) IOrderDetailDo
	Join(table schema.Tabler, on ...field.Expr) IOrderDetailDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IOrderDetailDo
	RightJoin(table schema.Tabler, on ...field.Expr) IOrderDetailDo
	Group(cols ...field.Expr) IOrderDetailDo
	Having(conds ...gen.Condition) IOrderDetailDo
	Limit(limit int) IOrderDetailDo
	Offset(offset int) IOrderDetailDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IOrderDetailDo
	Unscoped() IOrderDetailDo
	Create(values ...*model.OrderDetail) error
	CreateInBatches(values []*model.OrderDetail, batchSize int) error
	Save(values ...*model.OrderDetail) error
	First() (*model.OrderDetail, error)
	Take() (*model.OrderDetail, error)
	Last() (*model.OrderDetail, error)
	Find() ([]*model.OrderDetail, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.OrderDetail, err error)
	FindInBatches(result *[]*model.OrderDetail, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.OrderDetail) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IOrderDetailDo
	Assign(attrs ...field.AssignExpr) IOrderDetailDo
	Joins(fields ...field.RelationField) IOrderDetailDo
	Preload(fields ...field.RelationField) IOrderDetailDo
	FirstOrInit() (*model.OrderDetail, error)
	FirstOrCreate() (*model.OrderDetail, error)
	FindByPage(offset int, limit int) (result []*model.OrderDetail, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IOrderDetailDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (o orderDetailDo) Debug() IOrderDetailDo {
	return o.withDO(o.DO.Debug())
}

func (o orderDetailDo) WithContext(ctx context.Context) IOrderDetailDo {
	return o.withDO(o.DO.WithContext(ctx))
}

func (o orderDetailDo) ReadDB() IOrderDetailDo {
	return o.Clauses(dbresolver.Read)
}

func (o orderDetailDo) WriteDB() IOrderDetailDo {
	return o.Clauses(dbresolver.Write)
}

func (o orderDetailDo) Session(config *gorm.Session) IOrderDetailDo {
	return o.withDO(o.DO.Session(config))
}

func (o orderDetailDo) Clauses(conds ...clause.Expression) IOrderDetailDo {
	return o.withDO(o.DO.Clauses(conds...))
}

func (o orderDetailDo) Returning(value interface{}, columns ...string) IOrderDetailDo {
	return o.withDO(o.DO.Returning(value, columns...))
}

func (o orderDetailDo) Not(conds ...gen.Condition) IOrderDetailDo {
	return o.withDO(o.DO.Not(conds...))
}

func (o orderDetailDo) Or(conds ...gen.Condition) IOrderDetailDo {
	return o.withDO(o.DO.Or(conds...))
}

func (o orderDetailDo) Select(conds ...field.Expr) IOrderDetailDo {
	return o.withDO(o.DO.Select(conds...))
}

func (o orderDetailDo) Where(conds ...gen.Condition) IOrderDetailDo {
	return o.withDO(o.DO.Where(conds...))
}

func (o orderDetailDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) IOrderDetailDo {
	return o.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (o orderDetailDo) Order(conds ...field.Expr) IOrderDetailDo {
	return o.withDO(o.DO.Order(conds...))
}

func (o orderDetailDo) Distinct(cols ...field.Expr) IOrderDetailDo {
	return o.withDO(o.DO.Distinct(cols...))
}

func (o orderDetailDo) Omit(cols ...field.Expr) IOrderDetailDo {
	return o.withDO(o.DO.Omit(cols...))
}

func (o orderDetailDo) Join(table schema.Tabler, on ...field.Expr) IOrderDetailDo {
	return o.withDO(o.DO.Join(table, on...))
}

func (o orderDetailDo) LeftJoin(table schema.Tabler, on ...field.Expr) IOrderDetailDo {
	return o.withDO(o.DO.LeftJoin(table, on...))
}

func (o orderDetailDo) RightJoin(table schema.Tabler, on ...field.Expr) IOrderDetailDo {
	return o.withDO(o.DO.RightJoin(table, on...))
}

func (o orderDetailDo) Group(cols ...field.Expr) IOrderDetailDo {
	return o.withDO(o.DO.Group(cols...))
}

func (o orderDetailDo) Having(conds ...gen.Condition) IOrderDetailDo {
	return o.withDO(o.DO.Having(conds...))
}

func (o orderDetailDo) Limit(limit int) IOrderDetailDo {
	return o.withDO(o.DO.Limit(limit))
}

func (o orderDetailDo) Offset(offset int) IOrderDetailDo {
	return o.withDO(o.DO.Offset(offset))
}

func (o orderDetailDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IOrderDetailDo {
	return o.withDO(o.DO.Scopes(funcs...))
}

func (o orderDetailDo) Unscoped() IOrderDetailDo {
	return o.withDO(o.DO.Unscoped())
}

func (o orderDetailDo) Create(values ...*model.OrderDetail) error {
	if len(values) == 0 {
		return nil
	}
	return o.DO.Create(values)
}

func (o orderDetailDo) CreateInBatches(values []*model.OrderDetail, batchSize int) error {
	return o.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (o orderDetailDo) Save(values ...*model.OrderDetail) error {
	if len(values) == 0 {
		return nil
	}
	return o.DO.Save(values)
}

func (o orderDetailDo) First() (*model.OrderDetail, error) {
	if result, err := o.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.OrderDetail), nil
	}
}

func (o orderDetailDo) Take() (*model.OrderDetail, error) {
	if result, err := o.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.OrderDetail), nil
	}
}

func (o orderDetailDo) Last() (*model.OrderDetail, error) {
	if result, err := o.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.OrderDetail), nil
	}
}

func (o orderDetailDo) Find() ([]*model.OrderDetail, error) {
	result, err := o.DO.Find()
	return result.([]*model.OrderDetail), err
}

func (o orderDetailDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.OrderDetail, err error) {
	buf := make([]*model.OrderDetail, 0, batchSize)
	err = o.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (o orderDetailDo) FindInBatches(result *[]*model.OrderDetail, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return o.DO.FindInBatches(result, batchSize, fc)
}

func (o orderDetailDo) Attrs(attrs ...field.AssignExpr) IOrderDetailDo {
	return o.withDO(o.DO.Attrs(attrs...))
}

func (o orderDetailDo) Assign(attrs ...field.AssignExpr) IOrderDetailDo {
	return o.withDO(o.DO.Assign(attrs...))
}

func (o orderDetailDo) Joins(fields ...field.RelationField) IOrderDetailDo {
	for _, _f := range fields {
		o = *o.withDO(o.DO.Joins(_f))
	}
	return &o
}

func (o orderDetailDo) Preload(fields ...field.RelationField) IOrderDetailDo {
	for _, _f := range fields {
		o = *o.withDO(o.DO.Preload(_f))
	}
	return &o
}

func (o orderDetailDo) FirstOrInit() (*model.OrderDetail, error) {
	if result, err := o.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.OrderDetail), nil
	}
}

func (o orderDetailDo) FirstOrCreate() (*model.OrderDetail, error) {
	if result, err := o.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.OrderDetail), nil
	}
}

func (o orderDetailDo) FindByPage(offset int, limit int) (result []*model.OrderDetail, count int64, err error) {
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

func (o orderDetailDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = o.Count()
	if err != nil {
		return
	}

	err = o.Offset(offset).Limit(limit).Scan(result)
	return
}

func (o orderDetailDo) Scan(result interface{}) (err error) {
	return o.DO.Scan(result)
}

func (o orderDetailDo) Delete(models ...*model.OrderDetail) (result gen.ResultInfo, err error) {
	return o.DO.Delete(models)
}

func (o *orderDetailDo) withDO(do gen.Dao) *orderDetailDo {
	o.DO = *do.(*gen.DO)
	return o
}
