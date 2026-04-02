// Package trade 提供交易客户端，封装所有交易相关 API。
package trade

import (
	"encoding/json"

	"github.com/tigerfintech/openapi-go-sdk/client"
	"github.com/tigerfintech/openapi-go-sdk/model"
)

// TradeClient 交易客户端，封装所有交易相关 API。
type TradeClient struct {
	httpClient *client.HttpClient
	account    string
}

// NewTradeClient 创建交易客户端
func NewTradeClient(httpClient *client.HttpClient, account string) *TradeClient {
	return &TradeClient{httpClient: httpClient, account: account}
}

// execute 内部通用方法：构造请求、发送、返回 data 字段
func (c *TradeClient) execute(method string, bizParams interface{}) (json.RawMessage, error) {
	req, err := client.NewApiRequest(method, bizParams)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Execute(req)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// === 合约查询方法 ===

// Contract 查询单个合约
func (c *TradeClient) Contract(symbol, secType string) (json.RawMessage, error) {
	params := map[string]interface{}{
		"account": c.account,
		"symbol":  symbol,
		"secType": secType,
	}
	return c.execute("contract", params)
}

// Contracts 批量查询合约
func (c *TradeClient) Contracts(symbols []string, secType string) (json.RawMessage, error) {
	params := map[string]interface{}{
		"account": c.account,
		"symbols": symbols,
		"secType": secType,
	}
	return c.execute("contracts", params)
}

// QuoteContract 查询衍生品合约
func (c *TradeClient) QuoteContract(symbol, secType string) (json.RawMessage, error) {
	params := map[string]interface{}{
		"account": c.account,
		"symbol":  symbol,
		"secType": secType,
	}
	return c.execute("quote_contract", params)
}

// === 订单操作方法 ===

// PlaceOrder 下单
func (c *TradeClient) PlaceOrder(order model.Order) (json.RawMessage, error) {
	order.Account = c.account
	return c.execute("place_order", order)
}

// PreviewOrder 预览订单
func (c *TradeClient) PreviewOrder(order model.Order) (json.RawMessage, error) {
	order.Account = c.account
	return c.execute("preview_order", order)
}

// ModifyOrder 修改订单
func (c *TradeClient) ModifyOrder(id int64, order model.Order) (json.RawMessage, error) {
	order.Account = c.account
	order.ID = id
	return c.execute("modify_order", order)
}

// CancelOrder 取消订单
func (c *TradeClient) CancelOrder(id int64) (json.RawMessage, error) {
	params := map[string]interface{}{
		"account": c.account,
		"id":      id,
	}
	return c.execute("cancel_order", params)
}

// === 订单查询方法 ===

// Orders 查询全部订单
func (c *TradeClient) Orders() (json.RawMessage, error) {
	params := map[string]interface{}{"account": c.account}
	return c.execute("orders", params)
}

// ActiveOrders 查询待成交订单
func (c *TradeClient) ActiveOrders() (json.RawMessage, error) {
	params := map[string]interface{}{"account": c.account}
	return c.execute("active_orders", params)
}

// InactiveOrders 查询已撤销订单
func (c *TradeClient) InactiveOrders() (json.RawMessage, error) {
	params := map[string]interface{}{"account": c.account}
	return c.execute("inactive_orders", params)
}

// FilledOrders 查询已成交订单
func (c *TradeClient) FilledOrders() (json.RawMessage, error) {
	params := map[string]interface{}{"account": c.account}
	return c.execute("filled_orders", params)
}

// === 持仓和资产查询方法 ===

// Positions 查询持仓
func (c *TradeClient) Positions() (json.RawMessage, error) {
	params := map[string]interface{}{"account": c.account}
	return c.execute("positions", params)
}

// Assets 查询资产
func (c *TradeClient) Assets() (json.RawMessage, error) {
	params := map[string]interface{}{"account": c.account}
	return c.execute("assets", params)
}

// PrimeAssets 查询综合账户资产
func (c *TradeClient) PrimeAssets() (json.RawMessage, error) {
	params := map[string]interface{}{"account": c.account}
	return c.execute("prime_assets", params)
}

// OrderTransactions 查询订单成交明细
func (c *TradeClient) OrderTransactions(id int64) (json.RawMessage, error) {
	params := map[string]interface{}{
		"account": c.account,
		"id":      id,
	}
	return c.execute("order_transactions", params)
}
