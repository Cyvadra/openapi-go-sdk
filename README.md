# Tiger OpenAPI Go SDK

老虎证券 OpenAPI 的 Go 语言 SDK，提供行情查询、交易下单、账户管理和实时推送等功能。

[![Go Reference](https://pkg.go.dev/badge/github.com/tigerfintech/openapi-go-sdk.svg)](https://pkg.go.dev/github.com/tigerfintech/openapi-go-sdk)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## 安装

```bash
go get github.com/tigerfintech/openapi-go-sdk
```

要求 Go 1.20 或更高版本。

## Quick Start

以下是一个完整的示例，从配置到查询行情：

```go
package main

import (
	"fmt"
	"log"

	"github.com/tigerfintech/openapi-go-sdk/client"
	"github.com/tigerfintech/openapi-go-sdk/config"
	"github.com/tigerfintech/openapi-go-sdk/quote"
)

func main() {
	// 1. 创建配置（从 properties 文件加载）
	cfg, err := config.NewClientConfig(
		config.WithPropertiesFile("tiger_openapi_config.properties"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// 2. 创建 HTTP 客户端
	httpClient := client.NewHttpClient(cfg)

	// 3. 创建行情客户端并查询
	qc := quote.NewQuoteClient(httpClient)
	states, err := qc.GetMarketState("US")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("美股市场状态:", string(states))
}
```

## 配置

SDK 支持三种配置方式，优先级：**环境变量 > 代码设置（含配置文件） > 默认值**。

### 方式一：代码直接设置

```go
cfg, err := config.NewClientConfig(
	config.WithTigerID("你的 tiger_id"),
	config.WithPrivateKey("你的 RSA 私钥"),
	config.WithAccount("你的交易账户"),
)
```

### 方式二：从 properties 配置文件加载

```go
cfg, err := config.NewClientConfig(
	config.WithPropertiesFile("tiger_openapi_config.properties"),
)
```

配置文件格式：

```properties
tiger_id=你的开发者ID
private_key=你的RSA私钥
account=你的交易账户
```

### 方式三：环境变量

```bash
export TIGEROPEN_TIGER_ID=你的开发者ID
export TIGEROPEN_PRIVATE_KEY=你的RSA私钥
export TIGEROPEN_ACCOUNT=你的交易账户
```

### 配置项说明

| 配置项 | 说明 | 必填 | 默认值 |
|--------|------|------|--------|
| tiger_id | 开发者 ID | 是 | - |
| private_key | RSA 私钥 | 是 | - |
| account | 交易账户 | 否 | - |
| language | 语言（zh_CN/zh_TW/en_US） | 否 | zh_CN |
| timeout | 请求超时 | 否 | 15s |
| sandbox_debug | 是否使用沙箱环境 | 否 | false |

## 行情查询

```go
httpClient := client.NewHttpClient(cfg)
qc := quote.NewQuoteClient(httpClient)

// 获取市场状态
states, err := qc.GetMarketState("US")

// 获取实时报价
briefs, err := qc.GetBrief([]string{"AAPL", "TSLA"})

// 获取 K 线数据
klines, err := qc.GetKline("AAPL", "day")

// 获取分时数据
timeline, err := qc.GetTimeline([]string{"AAPL"})

// 获取深度行情
depth, err := qc.GetQuoteDepth("AAPL")

// 获取期权到期日
expiry, err := qc.GetOptionExpiration("AAPL")

// 获取期权链
chain, err := qc.GetOptionChain("AAPL", "2024-01-19")

// 获取期货交易所列表
exchanges, err := qc.GetFutureExchange()
```

## 交易操作

```go
tc := trade.NewTradeClient(httpClient, cfg.Account)

// 构造限价单
order := model.LimitOrder(cfg.Account, "AAPL", "STK", "BUY", 100, 150.0)

// 下单
result, err := tc.PlaceOrder(order)

// 预览订单（不实际下单）
preview, err := tc.PreviewOrder(order)

// 修改订单
modifiedOrder := order
modifiedOrder.LimitPrice = 155.0
result, err = tc.ModifyOrder(orderId, modifiedOrder)

// 取消订单
result, err = tc.CancelOrder(orderId)

// 查询全部订单
orders, err := tc.GetOrders()

// 查询持仓
positions, err := tc.GetPositions()

// 查询资产
assets, err := tc.GetAssets()
```

## 通用方法（ExecuteRaw）

当 SDK 尚未封装某个 API 时，可以使用 `ExecuteRaw` 直接调用：

```go
httpClient := client.NewHttpClient(cfg)

// 直接传入 API 方法名和 JSON 参数
resp, err := httpClient.ExecuteRaw("market_state", `{"market":"US"}`)
if err != nil {
	log.Fatal(err)
}
fmt.Println("原始响应:", resp)
```

## 实时推送

通过 WebSocket 长连接接收实时行情和账户推送，支持自动重连和心跳保活：

```go
import "github.com/tigerfintech/openapi-go-sdk/push"

pc := push.NewPushClient(cfg)

// 设置回调
pc.SetCallbacks(push.Callbacks{
	OnQuote: func(data *push.QuoteData) {
		fmt.Printf("行情推送: %s 最新价: %.2f\n", data.Symbol, data.LatestPrice)
	},
	OnOrder: func(data *push.OrderData) {
		fmt.Printf("订单推送: %s 状态: %s\n", data.Symbol, data.Status)
	},
	OnAsset: func(data *push.AssetData) {
		fmt.Println("资产变动推送")
	},
	OnPosition: func(data *push.PositionData) {
		fmt.Println("持仓变动推送")
	},
	OnConnect: func() {
		fmt.Println("推送连接成功")
	},
	OnDisconnect: func() {
		fmt.Println("推送连接断开")
	},
	OnError: func(err error) {
		fmt.Println("推送错误:", err)
	},
})

// 连接
if err := pc.Connect(); err != nil {
	log.Fatal(err)
}
defer pc.Disconnect()

// 订阅行情
pc.SubscribeQuote([]string{"AAPL", "TSLA"})

// 订阅账户推送
pc.SubscribeAsset("")       // 传空字符串使用配置中的账户
pc.SubscribeOrder("")
pc.SubscribePosition("")
```

## 项目结构

```
openapi-go-sdk/
├── config/    # 配置管理（ClientConfig、ConfigParser、动态域名）
├── signer/    # RSA 签名
├── client/    # HTTP 客户端（请求/响应、重试策略、ExecuteRaw）
├── model/     # 数据模型（Order、Contract、Position、枚举）
├── quote/     # 行情查询客户端
├── trade/     # 交易客户端
├── push/      # WebSocket 实时推送客户端
├── logger/    # 日志模块
└── examples/  # 示例代码
```

## API 参考

- [老虎证券 OpenAPI 文档](https://quant.itigerup.com/openapi/zh/python/overview/introduction.html)
- [Go SDK pkg.go.dev](https://pkg.go.dev/github.com/tigerfintech/openapi-go-sdk)

## 许可证

[MIT License](LICENSE)
