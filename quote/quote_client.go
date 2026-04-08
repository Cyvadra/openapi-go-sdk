// Package quote 提供行情查询客户端，封装所有行情相关 API。
package quote

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/tigerfintech/openapi-go-sdk/client"
)

const (
	apiVersionV2      = "2.0"
	apiVersionV3      = "3.0"
	defaultBarBeginMS = int64(-1)
	defaultBarEndMS   = int64(4070880000000)
)

var easternTZ = mustLoadLocation("America/New_York")

type strikePrice float64

func (value strikePrice) MarshalJSON() ([]byte, error) {
	formatted := strconv.FormatFloat(float64(value), 'f', -1, 64)
	if !strings.ContainsAny(formatted, ".eE") {
		formatted += ".0"
	}
	return []byte(formatted), nil
}

type optionBasicRequest struct {
	Symbol string      `json:"symbol"`
	Expiry int64       `json:"expiry"`
	Right  string      `json:"right,omitempty"`
	Strike strikePrice `json:"strike,omitempty"`
}

type optionKlineRequest struct {
	Symbol    string      `json:"symbol"`
	Expiry    int64       `json:"expiry"`
	Right     string      `json:"right,omitempty"`
	Strike    strikePrice `json:"strike,omitempty"`
	Period    string      `json:"period,omitempty"`
	BeginTime int64       `json:"begin_time,omitempty"`
	EndTime   int64       `json:"end_time,omitempty"`
}

// QuoteClient 行情查询客户端，封装所有行情相关 API。
type QuoteClient struct {
	httpClient  *client.HttpClient
	permissions json.RawMessage
}

// NewQuoteClient 创建行情查询客户端
func NewQuoteClient(httpClient *client.HttpClient) *QuoteClient {
	return &QuoteClient{httpClient: httpClient}
}

// NewQuoteClientWithPermissions 创建行情客户端并主动拉取行情权限。
func NewQuoteClientWithPermissions(httpClient *client.HttpClient) (*QuoteClient, error) {
	quoteClient := NewQuoteClient(httpClient)
	if err := quoteClient.InitializePermissions(); err != nil {
		return nil, err
	}
	return quoteClient, nil
}

// InitializePermissions 模拟 Python SDK 的初始化流程，先获取行情权限。
func (c *QuoteClient) InitializePermissions() error {
	if len(c.permissions) > 0 {
		return nil
	}
	permissions, err := c.GrabQuotePermission()
	if err != nil {
		return err
	}
	c.permissions = permissions
	return nil
}

// execute 内部通用方法：构造请求、发送、返回 data 字段
func (c *QuoteClient) execute(method string, bizParams interface{}) (json.RawMessage, error) {
	return c.executeWithVersion(method, bizParams, "")
}

func (c *QuoteClient) executeWithVersion(method string, bizParams interface{}, version string) (json.RawMessage, error) {
	req, err := client.NewVersionedApiRequest(method, bizParams, version)
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
	expiryMS, err := parseOptionExpiry(expiry)
	if err != nil {
		return nil, err
	}
	params := map[string]interface{}{
		"option_basic": []map[string]interface{}{{
			"symbol": symbol,
			"expiry": expiryMS,
		}},
		"return_greek_value": false,
	}
	return c.executeWithVersion("option_chain", params, apiVersionV3)
}

// OptionBrief 获取期权报价
func (c *QuoteClient) OptionBrief(identifiers []string) (json.RawMessage, error) {
	contracts := make([]optionBasicRequest, 0, len(identifiers))
	for _, identifier := range identifiers {
		contract, err := optionContractFromIdentifier(identifier)
		if err != nil {
			return nil, err
		}
		contracts = append(contracts, contract)
	}
	params := map[string]interface{}{"option_basic": contracts}
	return c.executeWithVersion("option_brief", params, apiVersionV2)
}

// OptionKline 获取期权 K 线
func (c *QuoteClient) OptionKline(identifier, period string) (json.RawMessage, error) {
	contract, err := optionContractFromIdentifier(identifier)
	if err != nil {
		return nil, err
	}
	params := map[string]interface{}{
		"option_query": []optionKlineRequest{{
			Symbol:    contract.Symbol,
			Expiry:    contract.Expiry,
			Right:     contract.Right,
			Strike:    contract.Strike,
			Period:    period,
			BeginTime: defaultBarBeginMS,
			EndTime:   defaultBarEndMS,
		}},
	}
	return c.executeWithVersion("option_kline", params, apiVersionV2)
}

func mustLoadLocation(name string) *time.Location {
	location, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return location
}

func parseOptionExpiry(expiry string) (int64, error) {
	expiry = strings.TrimSpace(expiry)
	if expiry == "" {
		return 0, fmt.Errorf("expiry is required")
	}
	if timestamp, err := strconv.ParseInt(expiry, 10, 64); err == nil {
		return timestamp, nil
	}
	parsed, err := time.ParseInLocation("2006-01-02", expiry, easternTZ)
	if err != nil {
		return 0, fmt.Errorf("invalid option expiry %q: %w", expiry, err)
	}
	return parsed.UnixMilli(), nil
}

func optionContractFromIdentifier(identifier string) (optionBasicRequest, error) {
	parts := strings.Fields(strings.TrimSpace(identifier))
	if len(parts) != 2 {
		return optionBasicRequest{}, fmt.Errorf("invalid option identifier %q", identifier)
	}
	contractCode := parts[1]
	if len(contractCode) != 15 {
		return optionBasicRequest{}, fmt.Errorf("invalid option identifier %q", identifier)
	}

	expiry, err := time.ParseInLocation("060102", contractCode[:6], easternTZ)
	if err != nil {
		return optionBasicRequest{}, fmt.Errorf("invalid option identifier %q: %w", identifier, err)
	}

	right := ""
	switch contractCode[6] {
	case 'C':
		right = "CALL"
	case 'P':
		right = "PUT"
	default:
		return optionBasicRequest{}, fmt.Errorf("invalid option identifier %q", identifier)
	}

	strikeInt, err := strconv.ParseInt(contractCode[7:], 10, 64)
	if err != nil {
		return optionBasicRequest{}, fmt.Errorf("invalid option identifier %q: %w", identifier, err)
	}

	return optionBasicRequest{
		Symbol: parts[0],
		Expiry: expiry.UnixMilli(),
		Right:  right,
		Strike: strikePrice(float64(strikeInt) / 1000),
	}, nil
}

func (c *QuoteClient) GetMarketState(market string) (json.RawMessage, error) {
	return c.MarketState(market)
}
func (c *QuoteClient) GetBrief(symbols []string) (json.RawMessage, error) {
	return c.QuoteRealTime(symbols)
}
func (c *QuoteClient) GetKline(symbol, period string) (json.RawMessage, error) {
	return c.Kline(symbol, period)
}
func (c *QuoteClient) GetTimeline(symbols []string) (json.RawMessage, error) {
	return c.Timeline(symbols)
}
func (c *QuoteClient) GetTradeTick(symbols []string) (json.RawMessage, error) {
	return c.TradeTick(symbols)
}
func (c *QuoteClient) GetQuoteDepth(symbol string) (json.RawMessage, error) {
	return c.QuoteDepth(symbol)
}
func (c *QuoteClient) GetOptionExpiration(symbol string) (json.RawMessage, error) {
	return c.OptionExpiration(symbol)
}
func (c *QuoteClient) GetOptionChain(symbol, expiry string) (json.RawMessage, error) {
	return c.OptionChain(symbol, expiry)
}
func (c *QuoteClient) GetOptionBrief(identifiers []string) (json.RawMessage, error) {
	return c.OptionBrief(identifiers)
}
func (c *QuoteClient) GetOptionKline(identifier, period string) (json.RawMessage, error) {
	return c.OptionKline(identifier, period)
}
func (c *QuoteClient) GetFutureExchange() (json.RawMessage, error) { return c.FutureExchange() }
func (c *QuoteClient) GetFutureContracts(exchange string) (json.RawMessage, error) {
	return c.FutureContracts(exchange)
}
func (c *QuoteClient) GetFutureRealTimeQuote(symbols []string) (json.RawMessage, error) {
	return c.FutureRealTimeQuote(symbols)
}
func (c *QuoteClient) GetFutureKline(symbol, period string) (json.RawMessage, error) {
	return c.FutureKline(symbol, period)
}
func (c *QuoteClient) GetFinancialDaily(symbol string) (json.RawMessage, error) {
	return c.FinancialDaily(symbol)
}
func (c *QuoteClient) GetFinancialReport(symbol string) (json.RawMessage, error) {
	return c.FinancialReport(symbol)
}
func (c *QuoteClient) GetCorporateAction(symbol string) (json.RawMessage, error) {
	return c.CorporateAction(symbol)
}
func (c *QuoteClient) GetCapitalFlow(symbol string) (json.RawMessage, error) {
	return c.CapitalFlow(symbol)
}
func (c *QuoteClient) GetCapitalDistribution(symbol string) (json.RawMessage, error) {
	return c.CapitalDistribution(symbol)
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
