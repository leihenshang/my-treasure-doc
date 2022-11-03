package order

import (
	"context"
	"errors"
	"fastduck/treasure-doc/service/mall/data/model"
	"fastduck/treasure-doc/service/mall/data/query"
	"fastduck/treasure-doc/service/mall/data/request/common"
	orderReq "fastduck/treasure-doc/service/mall/data/request/order"
	orderResp "fastduck/treasure-doc/service/mall/data/response/order"
	"fastduck/treasure-doc/service/mall/global"
	goodsDao "fastduck/treasure-doc/service/mall/internal/dao/goods"
	orderDao "fastduck/treasure-doc/service/mall/internal/dao/order"
	utils_datetime "fastduck/treasure-doc/service/mall/utils/utils-datetime"
	"fmt"
)

func OrderList(ctx context.Context, f orderReq.FilterOrderList) (res *orderResp.OrderList, err error) {
	res = new(orderResp.OrderList)
	if f.UserId <= 0 {
		err = errors.New("用户id不能为空")
		return
	}

	filter := &orderDao.OrderListFilter{
		UserId:   f.UserId,
		DataSort: f.DataSort,
		Status:   f.Status,
	}
	pagination := &f.Pagination

	result, total, err := orderDao.OrderList(ctx, filter, pagination)
	if err != nil {
		global.ZapSugar.Errorf("[OrderList|orderDao.OrderList]get order list err:%+v,filter:%+v,pagination:%+v", err, *filter, pagination)
		err = errors.New("获取用户订单数据失败")
	}

	for _, v := range result {
		res.List = append(res.List, &orderResp.OrderEntity{
			ID:        v.ID,
			OrderNo:   v.OrderNo,
			UserID:    v.UserID,
			Amount:    v.Amount,
			Status:    v.Status,
			CreatedAt: utils_datetime.TimeToFormat(&v.CreatedAt),
			UpdatedAt: utils_datetime.TimeToFormat(&v.UpdatedAt),
			DeletedAt: utils_datetime.TimeToFormat(&v.DeletedAt.Time),
		})
	}

	res.Total = total

	return
}

func OrderDetail(ctx context.Context, f orderReq.FilterOrderDetail) (res *orderResp.OrderInfo, err error) {
	if f.OrderId == 0 {
		err = errors.New("订单id不能为空")
		return
	}
	if f.UserId == 0 {
		err = errors.New("用户id不能为空")
		return
	}

	orderF := &orderDao.GetOrderFilter{
		OrderId: f.OrderId,
		UserId:  f.UserId,
	}

	result, err := orderDao.GetOrder(ctx, orderF)
	if err != nil {
		err = errors.New("获取订单失败")
		global.ZapSugar.Errorf("[OrderDetail|orderDao.GetOrder] failed to get order data.params:%+v,err:%+v", orderF, err)
		return
	}

	res = new(orderResp.OrderInfo)
	res.OrderEntity = &orderResp.OrderEntity{
		ID:        result.ID,
		OrderNo:   result.OrderNo,
		UserID:    result.UserID,
		Amount:    result.Amount,
		Status:    result.Status,
		CreatedAt: utils_datetime.TimeToFormat(&result.CreatedAt),
		UpdatedAt: utils_datetime.TimeToFormat(&result.UpdatedAt),
		DeletedAt: utils_datetime.TimeToFormat(&result.DeletedAt.Time),
	}

	// 获取订单明细
	details, _, err := orderDao.OrderDetailList(ctx, &orderDao.OrderDetailListFilter{
		OrderId: result.ID,
	}, &common.Pagination{
		Limit:  10000,
		Offset: 0,
	})
	if err != nil {
		err = errors.New("获取订单明细失败")
		global.ZapSugar.Errorf("[OrderDetail|orderDao.OrderDetailList] failed to get order data.params:%+v,err:%+v", orderF, err)
		return
	}
	for _, v := range details {
		res.OrderDetailEntity = append(res.OrderDetailEntity, &orderResp.OrderDetailEntity{
			ID:        v.ID,
			OrderID:   v.OrderID,
			GoodID:    v.GoodID,
			SkuID:     v.SkuID,
			Price:     v.Price,
			Quantity:  v.Quantity,
			CreatedAt: utils_datetime.TimeToFormat(&v.CreatedAt),
			UpdatedAt: utils_datetime.TimeToFormat(&v.UpdatedAt),
			DeletedAt: utils_datetime.TimeToFormat(&v.DeletedAt.Time),
		})
	}

	return
}

func OrderCreate(ctx context.Context, f orderReq.ParamsOrderCreate) (res *orderResp.OrderCreate, err error) {
	res = new(orderResp.OrderCreate)
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
		global.ZapSugar.Errorf("OrderCreate failed to get sku info. params:%+v err:%+v", *skuF, skuErr)
		err = errors.New("获取sku信息失败")
		return
	}
	if sku == nil {
		global.ZapSugar.Error("OrderCreate failed to get sku info. err:sku not existed")
		err = errors.New(fmt.Sprintf("sku不存在,skuId:%d", f.SkuId))
		return
	}

	remnant := sku.Stock - f.Quantity

	if remnant <= 0 {
		global.ZapSugar.Errorf("OrderCreate skuId:%d stock shortage", sku.ID)
		err = errors.New(fmt.Sprintf("库存不足,skuId:%d", f.SkuId))
		return
	}

	orderNo, noErr := GenerateOrderNo()
	if noErr != nil {
		global.ZapSugar.Errorf("[OrderCreate] generate order no err:%+v", noErr)
		err = errors.New("生成订单号错误")
		return
	}

	// 开启事务
	q := query.Use(global.DbIns)
	tx := q.Begin()

	_, skuUpdateErr := tx.WithContext(ctx).GoodsSku.Where(tx.GoodsSku.ID.Eq(sku.ID)).UpdateColumn(tx.GoodsSku.Stock, remnant)
	if skuUpdateErr != nil {
		err = errors.New("更新库存失败")
		global.ZapSugar.Error("[OrderCreate] update sku stock err,skuId:%+v,remnant:%+v", sku.ID, remnant)
		tx.Rollback()
		return
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
		global.ZapSugar.Error("[OrderCreate] create order err:%+v,data:%+v", orderErr, *insertOrder)
		tx.Rollback()
		return
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
		global.ZapSugar.Error("[OrderCreate] create order detail err:%+v,data:%#v", orderDetailErr, *insertOrderDetail)
		tx.Rollback()
		return
	}

	tx.Commit()

	res.OrderId = insertOrder.ID
	res.OrderNo = insertOrder.OrderNo

	return
}
