package order

import (
	"context"
	"errors"
	"fastduck/treasure-doc/service/mall/data/model"
	"fastduck/treasure-doc/service/mall/data/query"
	reqCommon "fastduck/treasure-doc/service/mall/data/request/common"
)

type OrderDetailListFilter struct {
	OrderId int32
	reqCommon.DataSort
}

func OrderDetailList(ctx context.Context, f *OrderDetailListFilter, p *reqCommon.Pagination) (result []*model.OrderDetail, total int64, err error) {
	if f == nil {
		err = errors.New("过滤器设置错误,不能为nil")
		return
	}

	if p == nil {
		p = reqCommon.NewPagination()
	}

	q := query.OrderDetail.WithContext(ctx)
	if f.OrderId > 0 {
		q = q.Where(query.OrderDetail.OrderID.Eq(f.OrderId))
	}

	if f.SortField != "" {
		orderCol, ok := query.OrderDetail.GetFieldByName(f.SortField)
		if ok {
			if f.IsDesc {
				q = q.Order(orderCol.Desc())
			} else {
				q = q.Order(orderCol)
			}
		}
	}

	result, total, err = q.FindByPage(int(p.Offset), int(p.Limit))
	return
}
