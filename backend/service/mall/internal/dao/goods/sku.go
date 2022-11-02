package goods

import (
	"context"
	"fastduck/treasure-doc/service/mall/data/model"
	"fastduck/treasure-doc/service/mall/data/query"
)

type GetSkuFilter struct {
	SkuId   int32
	Enabled int32
}

func GetSku(ctx context.Context, f *GetSkuFilter) (res *model.GoodsSku, err error) {
	q := query.GoodsSku.WithContext(ctx)
	if f.SkuId > 0 {
		q = q.Where(query.GoodsSku.ID.Eq(f.SkuId))
	}
	if f.Enabled > 0 {
		q = q.Where(query.GoodsSku.Enabled.Eq(f.Enabled))
	}

	res, err = q.First()

	return
}
