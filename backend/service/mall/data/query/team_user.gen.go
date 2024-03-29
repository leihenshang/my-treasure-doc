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

func newTeamUser(db *gorm.DB, opts ...gen.DOOption) teamUser {
	_teamUser := teamUser{}

	_teamUser.teamUserDo.UseDB(db, opts...)
	_teamUser.teamUserDo.UseModel(&model.TeamUser{})

	tableName := _teamUser.teamUserDo.TableName()
	_teamUser.ALL = field.NewAsterisk(tableName)
	_teamUser.ID = field.NewInt64(tableName, "id")
	_teamUser.UserID = field.NewInt64(tableName, "user_id")
	_teamUser.TeamID = field.NewInt64(tableName, "team_id")
	_teamUser.CreatedAt = field.NewTime(tableName, "created_at")
	_teamUser.UpdatedAt = field.NewTime(tableName, "updated_at")
	_teamUser.DeletedAt = field.NewField(tableName, "deleted_at")

	_teamUser.fillFieldMap()

	return _teamUser
}

type teamUser struct {
	teamUserDo

	ALL       field.Asterisk
	ID        field.Int64
	UserID    field.Int64
	TeamID    field.Int64
	CreatedAt field.Time
	UpdatedAt field.Time
	DeletedAt field.Field

	fieldMap map[string]field.Expr
}

func (t teamUser) Table(newTableName string) *teamUser {
	t.teamUserDo.UseTable(newTableName)
	return t.updateTableName(newTableName)
}

func (t teamUser) As(alias string) *teamUser {
	t.teamUserDo.DO = *(t.teamUserDo.As(alias).(*gen.DO))
	return t.updateTableName(alias)
}

func (t *teamUser) updateTableName(table string) *teamUser {
	t.ALL = field.NewAsterisk(table)
	t.ID = field.NewInt64(table, "id")
	t.UserID = field.NewInt64(table, "user_id")
	t.TeamID = field.NewInt64(table, "team_id")
	t.CreatedAt = field.NewTime(table, "created_at")
	t.UpdatedAt = field.NewTime(table, "updated_at")
	t.DeletedAt = field.NewField(table, "deleted_at")

	t.fillFieldMap()

	return t
}

func (t *teamUser) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := t.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (t *teamUser) fillFieldMap() {
	t.fieldMap = make(map[string]field.Expr, 6)
	t.fieldMap["id"] = t.ID
	t.fieldMap["user_id"] = t.UserID
	t.fieldMap["team_id"] = t.TeamID
	t.fieldMap["created_at"] = t.CreatedAt
	t.fieldMap["updated_at"] = t.UpdatedAt
	t.fieldMap["deleted_at"] = t.DeletedAt
}

func (t teamUser) clone(db *gorm.DB) teamUser {
	t.teamUserDo.ReplaceConnPool(db.Statement.ConnPool)
	return t
}

func (t teamUser) replaceDB(db *gorm.DB) teamUser {
	t.teamUserDo.ReplaceDB(db)
	return t
}

type teamUserDo struct{ gen.DO }

type ITeamUserDo interface {
	gen.SubQuery
	Debug() ITeamUserDo
	WithContext(ctx context.Context) ITeamUserDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() ITeamUserDo
	WriteDB() ITeamUserDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) ITeamUserDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) ITeamUserDo
	Not(conds ...gen.Condition) ITeamUserDo
	Or(conds ...gen.Condition) ITeamUserDo
	Select(conds ...field.Expr) ITeamUserDo
	Where(conds ...gen.Condition) ITeamUserDo
	Order(conds ...field.Expr) ITeamUserDo
	Distinct(cols ...field.Expr) ITeamUserDo
	Omit(cols ...field.Expr) ITeamUserDo
	Join(table schema.Tabler, on ...field.Expr) ITeamUserDo
	LeftJoin(table schema.Tabler, on ...field.Expr) ITeamUserDo
	RightJoin(table schema.Tabler, on ...field.Expr) ITeamUserDo
	Group(cols ...field.Expr) ITeamUserDo
	Having(conds ...gen.Condition) ITeamUserDo
	Limit(limit int) ITeamUserDo
	Offset(offset int) ITeamUserDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) ITeamUserDo
	Unscoped() ITeamUserDo
	Create(values ...*model.TeamUser) error
	CreateInBatches(values []*model.TeamUser, batchSize int) error
	Save(values ...*model.TeamUser) error
	First() (*model.TeamUser, error)
	Take() (*model.TeamUser, error)
	Last() (*model.TeamUser, error)
	Find() ([]*model.TeamUser, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.TeamUser, err error)
	FindInBatches(result *[]*model.TeamUser, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.TeamUser) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) ITeamUserDo
	Assign(attrs ...field.AssignExpr) ITeamUserDo
	Joins(fields ...field.RelationField) ITeamUserDo
	Preload(fields ...field.RelationField) ITeamUserDo
	FirstOrInit() (*model.TeamUser, error)
	FirstOrCreate() (*model.TeamUser, error)
	FindByPage(offset int, limit int) (result []*model.TeamUser, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) ITeamUserDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (t teamUserDo) Debug() ITeamUserDo {
	return t.withDO(t.DO.Debug())
}

func (t teamUserDo) WithContext(ctx context.Context) ITeamUserDo {
	return t.withDO(t.DO.WithContext(ctx))
}

func (t teamUserDo) ReadDB() ITeamUserDo {
	return t.Clauses(dbresolver.Read)
}

func (t teamUserDo) WriteDB() ITeamUserDo {
	return t.Clauses(dbresolver.Write)
}

func (t teamUserDo) Session(config *gorm.Session) ITeamUserDo {
	return t.withDO(t.DO.Session(config))
}

func (t teamUserDo) Clauses(conds ...clause.Expression) ITeamUserDo {
	return t.withDO(t.DO.Clauses(conds...))
}

func (t teamUserDo) Returning(value interface{}, columns ...string) ITeamUserDo {
	return t.withDO(t.DO.Returning(value, columns...))
}

func (t teamUserDo) Not(conds ...gen.Condition) ITeamUserDo {
	return t.withDO(t.DO.Not(conds...))
}

func (t teamUserDo) Or(conds ...gen.Condition) ITeamUserDo {
	return t.withDO(t.DO.Or(conds...))
}

func (t teamUserDo) Select(conds ...field.Expr) ITeamUserDo {
	return t.withDO(t.DO.Select(conds...))
}

func (t teamUserDo) Where(conds ...gen.Condition) ITeamUserDo {
	return t.withDO(t.DO.Where(conds...))
}

func (t teamUserDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) ITeamUserDo {
	return t.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (t teamUserDo) Order(conds ...field.Expr) ITeamUserDo {
	return t.withDO(t.DO.Order(conds...))
}

func (t teamUserDo) Distinct(cols ...field.Expr) ITeamUserDo {
	return t.withDO(t.DO.Distinct(cols...))
}

func (t teamUserDo) Omit(cols ...field.Expr) ITeamUserDo {
	return t.withDO(t.DO.Omit(cols...))
}

func (t teamUserDo) Join(table schema.Tabler, on ...field.Expr) ITeamUserDo {
	return t.withDO(t.DO.Join(table, on...))
}

func (t teamUserDo) LeftJoin(table schema.Tabler, on ...field.Expr) ITeamUserDo {
	return t.withDO(t.DO.LeftJoin(table, on...))
}

func (t teamUserDo) RightJoin(table schema.Tabler, on ...field.Expr) ITeamUserDo {
	return t.withDO(t.DO.RightJoin(table, on...))
}

func (t teamUserDo) Group(cols ...field.Expr) ITeamUserDo {
	return t.withDO(t.DO.Group(cols...))
}

func (t teamUserDo) Having(conds ...gen.Condition) ITeamUserDo {
	return t.withDO(t.DO.Having(conds...))
}

func (t teamUserDo) Limit(limit int) ITeamUserDo {
	return t.withDO(t.DO.Limit(limit))
}

func (t teamUserDo) Offset(offset int) ITeamUserDo {
	return t.withDO(t.DO.Offset(offset))
}

func (t teamUserDo) Scopes(funcs ...func(gen.Dao) gen.Dao) ITeamUserDo {
	return t.withDO(t.DO.Scopes(funcs...))
}

func (t teamUserDo) Unscoped() ITeamUserDo {
	return t.withDO(t.DO.Unscoped())
}

func (t teamUserDo) Create(values ...*model.TeamUser) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Create(values)
}

func (t teamUserDo) CreateInBatches(values []*model.TeamUser, batchSize int) error {
	return t.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (t teamUserDo) Save(values ...*model.TeamUser) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Save(values)
}

func (t teamUserDo) First() (*model.TeamUser, error) {
	if result, err := t.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.TeamUser), nil
	}
}

func (t teamUserDo) Take() (*model.TeamUser, error) {
	if result, err := t.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.TeamUser), nil
	}
}

func (t teamUserDo) Last() (*model.TeamUser, error) {
	if result, err := t.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.TeamUser), nil
	}
}

func (t teamUserDo) Find() ([]*model.TeamUser, error) {
	result, err := t.DO.Find()
	return result.([]*model.TeamUser), err
}

func (t teamUserDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.TeamUser, err error) {
	buf := make([]*model.TeamUser, 0, batchSize)
	err = t.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (t teamUserDo) FindInBatches(result *[]*model.TeamUser, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return t.DO.FindInBatches(result, batchSize, fc)
}

func (t teamUserDo) Attrs(attrs ...field.AssignExpr) ITeamUserDo {
	return t.withDO(t.DO.Attrs(attrs...))
}

func (t teamUserDo) Assign(attrs ...field.AssignExpr) ITeamUserDo {
	return t.withDO(t.DO.Assign(attrs...))
}

func (t teamUserDo) Joins(fields ...field.RelationField) ITeamUserDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Joins(_f))
	}
	return &t
}

func (t teamUserDo) Preload(fields ...field.RelationField) ITeamUserDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Preload(_f))
	}
	return &t
}

func (t teamUserDo) FirstOrInit() (*model.TeamUser, error) {
	if result, err := t.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.TeamUser), nil
	}
}

func (t teamUserDo) FirstOrCreate() (*model.TeamUser, error) {
	if result, err := t.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.TeamUser), nil
	}
}

func (t teamUserDo) FindByPage(offset int, limit int) (result []*model.TeamUser, count int64, err error) {
	result, err = t.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = t.Offset(-1).Limit(-1).Count()
	return
}

func (t teamUserDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = t.Count()
	if err != nil {
		return
	}

	err = t.Offset(offset).Limit(limit).Scan(result)
	return
}

func (t teamUserDo) Scan(result interface{}) (err error) {
	return t.DO.Scan(result)
}

func (t teamUserDo) Delete(models ...*model.TeamUser) (result gen.ResultInfo, err error) {
	return t.DO.Delete(models)
}

func (t *teamUserDo) withDO(do gen.Dao) *teamUserDo {
	t.DO = *do.(*gen.DO)
	return t
}
