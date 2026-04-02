// 交易下单示例
//
// 演示如何使用 TradeClient 进行下单、查询订单和持仓。
package examples

/*
import (
	"fmt"
	"log"

	"github.com/tigerfintech/openapi-go-sdk/config"
	"github.com/tigerfintech/openapi-go-sdk/model"
	"github.com/tigerfintech/openapi-go-sdk/trade"
)

func TradeExample() {
	cfg, err := config.NewClientConfig(
		config.WithTigerID("你的 tiger_id"),
		config.WithPrivateKey("你的 RSA 私钥"),
		config.WithAccount("你的交易账户"),
	)
	if err != nil {
		log.Fatal(err)
	}

	tc := trade.NewTradeClient(cfg)

	// 创建限价单
	order := model.LimitOrder("AAPL", "BUY", 100, 150.0)
	fmt.Printf("订单: %+v\n", order)

	// 下单（需要真实账户）
	// result, err := tc.PlaceOrder(order)

	// 查询订单
	orders, err := tc.GetOrders()
	if err != nil {
		log.Printf("查询订单失败: %v", err)
	} else {
		fmt.Printf("订单列表: %+v\n", orders)
	}

	// 查询持仓
	positions, err := tc.GetPositions()
	if err != nil {
		log.Printf("查询持仓失败: %v", err)
	} else {
		fmt.Printf("持仓列表: %+v\n", positions)
	}
}
*/
