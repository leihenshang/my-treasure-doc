package goods

import (
	"context"
	"errors"
	"fastduck/treasure-doc/service/mall/data/query"
	goodsReq "fastduck/treasure-doc/service/mall/data/request/goods"
	goodsResp "fastduck/treasure-doc/service/mall/data/response/goods"
	"fastduck/treasure-doc/service/mall/global"
	goodsDao "fastduck/treasure-doc/service/mall/internal/dao/goods"
	"strconv"
	"strings"
)

func GoodsList(ctx context.Context, f goodsReq.FilterGoodsList) (res *goodsResp.GoodsList, err error) {
	res = new(goodsResp.GoodsList)
	filter := &goodsDao.GoodsListFilter{
		GoodsName: f.GoodsName,
		DataSort:  f.DataSort,
	}
	pagination := &f.Pagination

	result, total, qErr := goodsDao.GoodsList(ctx, filter, pagination)
	if qErr != nil {
		global.ZapSugar.Errorf("GoodsList|goodsDao.GoodsList failed to get goods list.filter:%+v,pagination:%+v,err:%+v", filter, pagination, qErr)
		err = errors.New("获取商品列表失败")
		return
	}

	for _, v := range result {
		res.List = append(res.List, &goodsResp.GoodsEntity{
			ID:        v.ID,
			Img:       v.Img,
			GoodsName: v.GoodsName,
			CreatedAt: v.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: v.UpdatedAt.Format("2006-01-02 15:04:05"),
			DeletedAt: v.DeletedAt.Time.Format("2006-01-02 15:04:05"),
		})
	}

	res.Total = total

	return
}

func GoodsDetail(ctx context.Context, f goodsReq.FilterGoodsDetail) (res *goodsResp.GoodsDetail, err error) {
	if f.GoodsId == 0 {
		err = errors.New("商品id不能为空")
		return
	}

	q := query.Good.WithContext(ctx)
	result, qErr := q.Where(query.Good.ID.Eq(f.GoodsId)).First()
	if qErr != nil {
		global.ZapSugar.Errorf("GoodsDetail failed to get goods info.err:%+v", qErr)
		err = errors.New("查询商品信息失败")
		return
	}

	res = &goodsResp.GoodsDetail{
		GoodsEntity: goodsResp.GoodsEntity{
			ID:        result.ID,
			Img:       result.Img,
			GoodsName: result.GoodsName,
			Quantity:  0,
			CreatedAt: result.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: result.UpdatedAt.Format("2006-01-02 15:04:05"),
			DeletedAt: result.DeletedAt.Time.Format("2006-01-02 15:04:05"),
		},
		Sku: []*goodsResp.GoodsSkuEntity{},
	}

	// sku 信息
	sku, skuErr := query.GoodsSku.WithContext(ctx).Where(query.GoodsSku.GoodsID.Eq(res.ID)).Find()
	if skuErr != nil {
		global.ZapSugar.Errorf("GoodsDetail failed to get goods sku.err:%+v", skuErr)
	}

	// 获取spec 描述
	spec, specErr := query.GoodsSpec.WithContext(ctx).Where(query.GoodsSpec.GoodID.In(res.ID)).Find()
	if specErr != nil {
		global.ZapSugar.Errorf("GoodsDetail failed to get goods spec.err:%+v", specErr)
	}

	specMap := make(map[int32][]string, 0)
	for _, v := range spec {
		specMap[v.ID] = []string{v.Spec, v.SpecVal, v.Units}
	}

	for _, v := range sku {
		res.Sku = append(res.Sku, &goodsResp.GoodsSkuEntity{
			ID:           v.ID,
			GoodsID:      v.GoodsID,
			GoodsSpecIds: v.GoodsSpecIds,
			GoodsSpec: func() (s []*goodsResp.SpecDesc) {
				goodsSpec := strings.Split(v.GoodsSpecIds, ",")
				for _, v := range goodsSpec {
					tempId, _ := strconv.Atoi(v)
					if vv, ok := specMap[int32(tempId)]; ok {
						s = append(s, &goodsResp.SpecDesc{
							SpecId:  int32(tempId),
							Spec:    vv[0],
							SpecVal: vv[1],
							Units:   vv[2],
						})
					}
				}

				return
			}(),
			Price:     v.Price,
			Stock:     v.Stock,
			CreatedAt: v.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: v.UpdatedAt.Format("2006-01-02 15:04:05"),
			DeletedAt: v.DeletedAt.Time.Format("2006-01-02 15:04:05"),
		})
	}

	return
}
