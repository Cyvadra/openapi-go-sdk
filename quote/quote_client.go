// Package quote 提供行情查询客户端，封装所有行情相关 API。
package quote

import (
	"encoding/json"

	"github.com/tigerfintech/openapi-go-sdk/client"
)

// QuoteClient 行情查询客户端，封装所有行情相关 API。
type QuoteClient struct {
	httpClient *client.HttpClient
}

// NewQuoteClient 创建行情查询客户端
func NewQuoteClient(httpClient *client.HttpClient) *QuoteClient {
	return &QuoteClient{httpClient: httpClient}
}

// execute 内部通用方法：构造请求、发送、返回 data 字段
func (c *QuoteClient) execute(method string, bizParams interface{}) (json.RawMessage, error) {
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

// === 基础行情方法 ===

// MarketState 获取市场状态
func (c *QuoteClient) MarketState(market string) (json.RawMessage, error) {
	params := map[string]interface{}{"market": market}
	return c.execute("market_state", params)
}

// QuoteRealTime 获取实时报价
func (c *QuoteClient) QuoteRealTime(symbols []string) (json.RawMessage, error) {
	params := map[string]interface{}{"symbols": symbols}
	return c.execute("quote_real_time", params)
}

// Kline 获取 K 线数据
func (c *QuoteClient) Kline(symbol, period string) (json.RawMessage, error) {
	params := map[string]interface{}{
		"symbols": []string{symbol},
		"period":  period,
	}
	return c.execute("kline", params)
}

// Timeline 获取分时数据
func (c *QuoteClient) Timeline(symbols []string) (json.RawMessage, error) {
	params := map[string]interface{}{"symbols": symbols}
	return c.execute("timeline", params)
}

// TradeTick 获取逐笔成交数据
func (c *QuoteClient) TradeTick(symbols []string) (json.RawMessage, error) {
	params := map[string]interface{}{"symbols": symbols}
	return c.execute("trade_tick", params)
}

// QuoteDepth 获取深度行情
func (c *QuoteClient) QuoteDepth(symbol string) (json.RawMessage, error) {
	params := map[string]interface{}{"symbol": symbol}
	return c.execute("quote_depth", params)
}

// === 期权行情方法 ===

// OptionExpiration 获取期权到期日
func (c *QuoteClient) OptionExpiration(symbol string) (json.RawMessage, error) {
	params := map[string]interface{}{"symbols": []string{symbol}}
	return c.execute("option_expiration", params)
}

// OptionChain 获取期权链
func (c *QuoteClient) OptionChain(symbol, expiry string) (json.RawMessage, error) {
	params := map[string]interface{}{
		"symbol": symbol,
		"expiry": expiry,
	}
	return c.execute("option_chain", params)
}

// OptionBrief 获取期权报价
func (c *QuoteClient) OptionBrief(identifiers []string) (json.RawMessage, error) {
	params := map[string]interface{}{"identifiers": identifiers}
	return c.execute("option_brief", params)
}

// OptionKline 获取期权 K 线
func (c *QuoteClient) OptionKline(identifier, period string) (json.RawMessage, error) {
	params := map[string]interface{}{
		"identifier": identifier,
		"period":     period,
	}
	return c.execute("option_kline", params)
}

// === 期货行情方法 ===

// FutureExchange 获取期货交易所列表
func (c *QuoteClient) FutureExchange() (json.RawMessage, error) {
	params := map[string]interface{}{"sec_type": "FUT"}
	return c.execute("future_exchange", params)
}

// FutureContracts 获取期货合约列表
func (c *QuoteClient) FutureContracts(exchange string) (json.RawMessage, error) {
	params := map[string]interface{}{"exchange": exchange}
	return c.execute("future_contracts", params)
}

// FutureRealTimeQuote 获取期货实时报价
func (c *QuoteClient) FutureRealTimeQuote(symbols []string) (json.RawMessage, error) {
	params := map[string]interface{}{"symbols": symbols}
	return c.execute("future_real_time_quote", params)
}

// FutureKline 获取期货 K 线
func (c *QuoteClient) FutureKline(symbol, period string) (json.RawMessage, error) {
	params := map[string]interface{}{
		"symbol": symbol,
		"period": period,
	}
	return c.execute("future_kline", params)
}

// === 基本面和资金流向方法 ===

// FinancialDaily 获取财务日报
func (c *QuoteClient) FinancialDaily(symbol string) (json.RawMessage, error) {
	params := map[string]interface{}{"symbol": symbol}
	return c.execute("financial_daily", params)
}

// FinancialReport 获取财务报告
func (c *QuoteClient) FinancialReport(symbol string) (json.RawMessage, error) {
	params := map[string]interface{}{"symbol": symbol}
	return c.execute("financial_report", params)
}

// CorporateAction 获取公司行动
func (c *QuoteClient) CorporateAction(symbol string) (json.RawMessage, error) {
	params := map[string]interface{}{"symbol": symbol}
	return c.execute("corporate_action", params)
}

// CapitalFlow 获取资金流向
func (c *QuoteClient) CapitalFlow(symbol string) (json.RawMessage, error) {
	params := map[string]interface{}{"symbol": symbol}
	return c.execute("capital_flow", params)
}

// CapitalDistribution 获取资金分布
func (c *QuoteClient) CapitalDistribution(symbol string) (json.RawMessage, error) {
	params := map[string]interface{}{"symbol": symbol}
	return c.execute("capital_distribution", params)
}

// === 选股器和行情权限方法 ===

// MarketScanner 选股器
func (c *QuoteClient) MarketScanner(params map[string]interface{}) (json.RawMessage, error) {
	return c.execute("market_scanner", params)
}

// GrabQuotePermission 获取行情权限
func (c *QuoteClient) GrabQuotePermission() (json.RawMessage, error) {
	return c.execute("grab_quote_permission", nil)
}
