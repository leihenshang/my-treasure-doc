package order

import (
	"context"
	"errors"
	"fastduck/treasure-doc/service/mall/data/model"
	"fastduck/treasure-doc/service/mall/data/query"
	orderReq "fastduck/treasure-doc/service/mall/data/request/order"
	orderResp "fastduck/treasure-doc/service/mall/data/response/order"
	"fastduck/treasure-doc/service/mall/global"
	goodsDao "fastduck/treasure-doc/service/mall/internal/dao/goods"
	"fmt"
)

func OrderList(ctx context.Context, f orderReq.FilterOrderList) (res []*model.Order, total int64, err error) {
	q := query.Order.WithContext(ctx)

	res, total, err = q.FindByPage(int(f.Offset), int(f.Limit))
	return
}

func OrderDetail(ctx context.Context, f orderReq.FilterOrderDetail) (res *model.Order, err error) {
	if f.OrderId == 0 {
		err = errors.New("订单id不能为空")
		return
	}

	q := query.Order.WithContext(ctx)
	res, err = q.Where(query.Order.ID.Eq(f.OrderId)).First()
	return
}

func OrderCreate(ctx context.Context, f orderReq.FilterOrderCreate) (res *orderResp.OrderCreate, err error) {
	if f.Quantity <= 0 {
		err = errors.New("数量不能小于1")
		return
	}
	if f.SkuId <= 0 {
		err = errors.New("skuId 不能小于1")
		return
	}

	if f.UserId <= 0 {
		err = errors.New("userId 不能小于1")
		return
	}

	//查询商品 sku 以及库存
	//扣减库存
	//生成订单
	skuF := &goodsDao.GetSkuFilter{
		SkuId:   f.SkuId,
		Enabled: 1,
	}
	sku, skuErr := goodsDao.GetSku(ctx, skuF)
	if skuErr != nil {
		global.ZAPSUGAR.Errorf("OrderCreate failed to get sku info. params:%+v err:%+v", *skuF, skuErr)
		err = errors.New("获取sku信息失败")
		return
	}
	if sku == nil {
		global.ZAPSUGAR.Error("OrderCreate failed to get sku info. err:sku not existed")
		err = errors.New(fmt.Sprintf("sku不存在,skuId:%d", f.SkuId))
		return
	}

	remnant := sku.Stock - f.Quantity

	if remnant <= 0 {
		global.ZAPSUGAR.Errorf("OrderCreate skuId:%d stock shortage", sku.ID)
		err = errors.New(fmt.Sprintf("库存不足,skuId:%d", f.SkuId))
	}

	orderNo, noErr := GenerateOrderNo()
	if noErr != nil {
		global.ZAPSUGAR.Errorf("[OrderCreate] generate order no err:%+v", noErr)
		err = errors.New("生成订单号错误")
	}

	// 开启事务
	q := query.Use(global.DB)
	tx := q.Begin()

	_, skuUpdateErr := tx.WithContext(ctx).GoodsSku.Where(tx.GoodsSku.ID.Eq(sku.ID)).UpdateColumn(tx.GoodsSku.Stock, remnant)
	if skuUpdateErr != nil {
		err = errors.New("更新库存失败")
		global.ZAPSUGAR.Error("[OrderCreate] update sku stock err,skuId:%+v,remnant:%+v", sku.ID, remnant)
		tx.Rollback()
	}

	// 订单

	insertOrder := &model.Order{
		OrderNo: orderNo,
		UserID:  f.UserId,
		Amount:  sku.Price * float64(f.Quantity),
		Status:  1,
	}

	orderErr := tx.WithContext(ctx).Order.Create(insertOrder)
	if orderErr != nil {
		err = errors.New("保存订单失败")
		global.ZAPSUGAR.Error("[OrderCreate] create order err:%+v,data:%+v", orderErr, *insertOrder)
		tx.Rollback()
	}

	insertOrderDetail := &model.OrderDetail{
		OrderID:  insertOrder.ID,
		GoodID:   sku.GoodsID,
		SkuID:    sku.ID,
		Price:    sku.Price,
		Quantity: f.Quantity,
	}
	orderDetailErr := tx.WithContext(ctx).OrderDetail.Create(insertOrderDetail)
	if orderDetailErr != nil {
		err = errors.New("保存订单明细失败")
		global.ZAPSUGAR.Error("[OrderCreate] create order detail err:%+v,data:%+v", orderDetailErr, *insertOrderDetail)
		tx.Rollback()
	}

	tx.Commit()

	return
}
