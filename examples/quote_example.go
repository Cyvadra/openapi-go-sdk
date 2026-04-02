// 行情查询示例
//
// 演示如何使用 QuoteClient 查询市场状态、实时报价和 K 线数据。
package examples

/*
import (
	"fmt"
	"log"

	"github.com/tigerfintech/openapi-sdks/go/config"
	"github.com/tigerfintech/openapi-sdks/go/quote"
)

func QuoteExample() {
	// 创建配置
	cfg, err := config.NewClientConfig(
		config.WithTigerID("你的 tiger_id"),
		config.WithPrivateKey("你的 RSA 私钥"),
		config.WithAccount("你的交易账户"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// 创建行情客户端
	qc := quote.NewQuoteClient(cfg)

	// 查询市场状态
	fmt.Println("=== 市场状态 ===")
	states, err := qc.GetMarketState("US")
	if err != nil {
		log.Printf("查询市场状态失败: %v", err)
	} else {
		fmt.Printf("市场状态: %+v\n", states)
	}

	// 查询实时报价
	fmt.Println("\n=== 实时报价 ===")
	quotes, err := qc.GetQuoteRealTime([]string{"AAPL", "TSLA", "GOOG"})
	if err != nil {
		log.Printf("查询实时报价失败: %v", err)
	} else {
		fmt.Printf("报价数据: %+v\n", quotes)
	}

	// 查询 K 线
	fmt.Println("\n=== K 线数据 ===")
	klines, err := qc.GetKline("AAPL", "day", 10)
	if err != nil {
		log.Printf("查询 K 线失败: %v", err)
	} else {
		fmt.Printf("K 线数据: %+v\n", klines)
	}
}
*/
