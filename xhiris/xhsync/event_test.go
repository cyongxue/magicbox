package xhsync

import (
	"context"
	"testing"
)

func TestEngine_Dispatch(t *testing.T) {
	// 事件运行引擎初始化
	cancelCtx, cancel := context.WithCancel(context.Background())
	if err := NewGroupEngine(cancelCtx, cancel, 8, 1024); err != nil {
		panic(err)
	}

	//parentGoodsPriceEv := ParentGoodsPriceEvent{
	//	Params: ParentGoodsPriceParams{
	//		Ctx:     ctx,
	//		GoodsId: input.Goods.ParentId,
	//	},
	//}
	//if err := runevent.EnginePrt.Dispatch(input.Goods.ParentId, &parentGoodsPriceEv); err != nil {
	//	mylog.Warn(ctx, mylog.DFLAG_UNDEF, map[string]interface{}{
	//		mylog.Args:     input.Goods.ParentId,
	//		"event":        "ParentGoodsPriceEvent",
	//		mylog.ErrorMsg: err.Error(),
	//	})
	//}
}
