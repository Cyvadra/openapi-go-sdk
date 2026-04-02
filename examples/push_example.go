// 实时推送示例
//
// 演示如何使用 PushClient 接收实时行情和账户推送。
package examples

/*
import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tigerfintech/openapi-sdks/go/config"
	"github.com/tigerfintech/openapi-sdks/go/push"
)

func PushExample() {
	cfg, err := config.NewClientConfig(
		config.WithPropertiesFile("tiger_openapi_config.properties"),
	)
	if err != nil {
		log.Fatal(err)
	}

	pc := push.NewPushClient(cfg)
	pc.SetCallbacks(push.Callbacks{
		OnQuote: func(data *push.QuoteData) {
			fmt.Printf("[行情] %s: %.2f\n", data.Symbol, data.LatestPrice)
		},
		OnOrder: func(data *push.OrderData) {
			fmt.Printf("[订单] %s: %s\n", data.Symbol, data.Status)
		},
		OnConnect:    func() { fmt.Println("已连接推送服务器") },
		OnDisconnect: func() { fmt.Println("已断开连接") },
		OnError:      func(err error) { fmt.Printf("错误: %v\n", err) },
	})

	if err := pc.Connect(); err != nil {
		log.Fatal(err)
	}
	defer pc.Disconnect()

	pc.SubscribeQuote([]string{"AAPL", "TSLA"})
	pc.SubscribeAsset("")

	// 等待退出信号
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
*/
