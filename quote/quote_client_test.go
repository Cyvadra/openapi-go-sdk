package quote

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"time"

	"github.com/tigerfintech/openapi-go-sdk/client"
	"github.com/tigerfintech/openapi-go-sdk/config"
)

// mustGenerateKeyPEM 生成测试用 PKCS#1 PEM 格式私钥
func mustGenerateKeyPEM() string {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic("生成 RSA 密钥对失败: " + err.Error())
	}
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privDER,
	})
	return string(privPEM)
}

var cachedTestPrivateKey = mustGenerateKeyPEM()

// mockServer 创建一个 mock HTTP 服务器，返回指定的 data
func mockServer(t *testing.T, expectedMethod string, data interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dataBytes, _ := json.Marshal(data)
		resp := map[string]interface{}{
			"code":      0,
			"message":   "success",
			"data":      json.RawMessage(dataBytes),
			"timestamp": 1700000000,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
}

// mockErrorServer 创建一个返回错误的 mock HTTP 服务器
func mockErrorServer(t *testing.T, code int, message string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"code":      code,
			"message":   message,
			"timestamp": 1700000000,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
}

func newTestQuoteClient(serverURL string) *QuoteClient {
	cfg := &config.ClientConfig{
		TigerID:    "test_tiger_id",
		PrivateKey: cachedTestPrivateKey,
		Account:    "test_account",
		Language:   "zh_CN",
		Timeout:    5 * time.Second,
		ServerURL:  serverURL,
	}
	httpClient := client.NewHttpClient(cfg)
	return NewQuoteClient(httpClient)
}

// === 14.1 基础行情测试 ===

func TestGetMarketState(t *testing.T) {
	server := mockServer(t, "market_state", []map[string]interface{}{
		{"market": "US", "status": "Trading"},
	})
	defer server.Close()

	qc := newTestQuoteClient(server.URL)
	data, err := qc.GetMarketState("US")
	if err != nil {
		t.Fatalf("GetMarketState 失败: %v", err)
	}
	if data == nil {
		t.Fatal("GetMarketState 返回 nil data")
	}
}

func TestGetBrief(t *testing.T) {
	server := mockServer(t, "brief", []map[string]interface{}{
		{"symbol": "AAPL", "latestPrice": 150.0},
	})
	defer server.Close()

	qc := newTestQuoteClient(server.URL)
	data, err := qc.GetBrief([]string{"AAPL"})
	if err != nil {
		t.Fatalf("GetBrief 失败: %v", err)
	}
	if data == nil {
		t.Fatal("GetBrief 返回 nil data")
	}
}

func TestGetKline(t *testing.T) {
	server := mockServer(t, "kline", []map[string]interface{}{
		{"symbol": "AAPL", "period": "day"},
	})
	defer server.Close()

	qc := newTestQuoteClient(server.URL)
	data, err := qc.GetKline("AAPL", "day")
	if err != nil {
		t.Fatalf("GetKline 失败: %v", err)
	}
	if data == nil {
		t.Fatal("GetKline 返回 nil data")
	}
}

func TestGetTimeline(t *testing.T) {
	server := mockServer(t, "timeline", []map[string]interface{}{
		{"symbol": "AAPL"},
	})
	defer server.Close()

	qc := newTestQuoteClient(server.URL)
	data, err := qc.GetTimeline([]string{"AAPL"})
	if err != nil {
		t.Fatalf("GetTimeline 失败: %v", err)
	}
	if data == nil {
		t.Fatal("GetTimeline 返回 nil data")
	}
}

func TestGetTradeTick(t *testing.T) {
	server := mockServer(t, "trade_tick", []map[string]interface{}{
		{"symbol": "AAPL"},
	})
	defer server.Close()

	qc := newTestQuoteClient(server.URL)
	data, err := qc.GetTradeTick([]string{"AAPL"})
	if err != nil {
		t.Fatalf("GetTradeTick 失败: %v", err)
	}
	if data == nil {
		t.Fatal("GetTradeTick 返回 nil data")
	}
}

func TestGetQuoteDepth(t *testing.T) {
	server := mockServer(t, "quote_depth", map[string]interface{}{
		"symbol": "AAPL",
	})
	defer server.Close()

	qc := newTestQuoteClient(server.URL)
	data, err := qc.GetQuoteDepth("AAPL")
	if err != nil {
		t.Fatalf("GetQuoteDepth 失败: %v", err)
	}
	if data == nil {
		t.Fatal("GetQuoteDepth 返回 nil data")
	}
}

// === 14.3 期权行情测试 ===

func TestGetOptionExpiration(t *testing.T) {
	server := mockServer(t, "option_expiration", []string{"2024-01-19", "2024-02-16"})
	defer server.Close()

	qc := newTestQuoteClient(server.URL)
	data, err := qc.GetOptionExpiration("AAPL")
	if err != nil {
		t.Fatalf("GetOptionExpiration 失败: %v", err)
	}
	if data == nil {
		t.Fatal("GetOptionExpiration 返回 nil data")
	}
}

func TestGetOptionChain(t *testing.T) {
	server := mockServer(t, "option_chain", []map[string]interface{}{
		{"symbol": "AAPL", "expiry": "2024-01-19"},
	})
	defer server.Close()

	qc := newTestQuoteClient(server.URL)
	data, err := qc.GetOptionChain("AAPL", "2024-01-19")
	if err != nil {
		t.Fatalf("GetOptionChain 失败: %v", err)
	}
	if data == nil {
		t.Fatal("GetOptionChain 返回 nil data")
	}
}

func TestGetOptionBrief(t *testing.T) {
	server := mockServer(t, "option_brief", []map[string]interface{}{
		{"identifier": "AAPL 240119C00150000"},
	})
	defer server.Close()

	qc := newTestQuoteClient(server.URL)
	data, err := qc.GetOptionBrief([]string{"AAPL 240119C00150000"})
	if err != nil {
		t.Fatalf("GetOptionBrief 失败: %v", err)
	}
	if data == nil {
		t.Fatal("GetOptionBrief 返回 nil data")
	}
}

func TestGetOptionKline(t *testing.T) {
	server := mockServer(t, "option_kline", []map[string]interface{}{
		{"identifier": "AAPL 240119C00150000", "period": "day"},
	})
	defer server.Close()

	qc := newTestQuoteClient(server.URL)
	data, err := qc.GetOptionKline("AAPL 240119C00150000", "day")
	if err != nil {
		t.Fatalf("GetOptionKline 失败: %v", err)
	}
	if data == nil {
		t.Fatal("GetOptionKline 返回 nil data")
	}
}

func TestOptionRequestsPreserveDecimalStrike(t *testing.T) {
	bodies := make([]string, 0, 2)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read request body: %v", err)
		}
		bodies = append(bodies, string(body))
		resp := map[string]interface{}{
			"code":      0,
			"message":   "success",
			"data":      []map[string]interface{}{},
			"timestamp": 1700000000,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	qc := newTestQuoteClient(server.URL)
	if _, err := qc.GetOptionBrief([]string{"AAPL 260417C00200000"}); err != nil {
		t.Fatalf("GetOptionBrief failed: %v", err)
	}
	if _, err := qc.GetOptionKline("AAPL 260417C00200000", "day"); err != nil {
		t.Fatalf("GetOptionKline failed: %v", err)
	}

	if len(bodies) != 2 {
		t.Fatalf("expected 2 request bodies, got %d", len(bodies))
	}
	for _, body := range bodies {
		if !strings.Contains(body, `\"strike\":200.0`) {
			t.Fatalf("expected decimal strike in biz_content, got %s", body)
		}
	}
}

// === 14.5 期货行情测试 ===

func TestGetFutureExchange(t *testing.T) {
	server := mockServer(t, "future_exchange", []string{"CME", "NYMEX"})
	defer server.Close()

	qc := newTestQuoteClient(server.URL)
	data, err := qc.GetFutureExchange()
	if err != nil {
		t.Fatalf("GetFutureExchange 失败: %v", err)
	}
	if data == nil {
		t.Fatal("GetFutureExchange 返回 nil data")
	}
}

func TestGetFutureContracts(t *testing.T) {
	server := mockServer(t, "future_contracts", []map[string]interface{}{
		{"symbol": "ES", "exchange": "CME"},
	})
	defer server.Close()

	qc := newTestQuoteClient(server.URL)
	data, err := qc.GetFutureContracts("CME")
	if err != nil {
		t.Fatalf("GetFutureContracts 失败: %v", err)
	}
	if data == nil {
		t.Fatal("GetFutureContracts 返回 nil data")
	}
}

func TestGetFutureRealTimeQuote(t *testing.T) {
	server := mockServer(t, "future_real_time_quote", []map[string]interface{}{
		{"symbol": "ES2312"},
	})
	defer server.Close()

	qc := newTestQuoteClient(server.URL)
	data, err := qc.GetFutureRealTimeQuote([]string{"ES2312"})
	if err != nil {
		t.Fatalf("GetFutureRealTimeQuote 失败: %v", err)
	}
	if data == nil {
		t.Fatal("GetFutureRealTimeQuote 返回 nil data")
	}
}

func TestGetFutureKline(t *testing.T) {
	server := mockServer(t, "future_kline", []map[string]interface{}{
		{"symbol": "ES2312", "period": "day"},
	})
	defer server.Close()

	qc := newTestQuoteClient(server.URL)
	data, err := qc.GetFutureKline("ES2312", "day")
	if err != nil {
		t.Fatalf("GetFutureKline 失败: %v", err)
	}
	if data == nil {
		t.Fatal("GetFutureKline 返回 nil data")
	}
}

// === 14.7 基本面和资金流向测试 ===

func TestGetFinancialDaily(t *testing.T) {
	server := mockServer(t, "financial_daily", []map[string]interface{}{
		{"symbol": "AAPL"},
	})
	defer server.Close()

	qc := newTestQuoteClient(server.URL)
	data, err := qc.GetFinancialDaily("AAPL")
	if err != nil {
		t.Fatalf("GetFinancialDaily 失败: %v", err)
	}
	if data == nil {
		t.Fatal("GetFinancialDaily 返回 nil data")
	}
}

func TestGetFinancialReport(t *testing.T) {
	server := mockServer(t, "financial_report", []map[string]interface{}{
		{"symbol": "AAPL"},
	})
	defer server.Close()

	qc := newTestQuoteClient(server.URL)
	data, err := qc.GetFinancialReport("AAPL")
	if err != nil {
		t.Fatalf("GetFinancialReport 失败: %v", err)
	}
	if data == nil {
		t.Fatal("GetFinancialReport 返回 nil data")
	}
}

func TestGetCorporateAction(t *testing.T) {
	server := mockServer(t, "corporate_action", []map[string]interface{}{
		{"symbol": "AAPL"},
	})
	defer server.Close()

	qc := newTestQuoteClient(server.URL)
	data, err := qc.GetCorporateAction("AAPL")
	if err != nil {
		t.Fatalf("GetCorporateAction 失败: %v", err)
	}
	if data == nil {
		t.Fatal("GetCorporateAction 返回 nil data")
	}
}

func TestGetCapitalFlow(t *testing.T) {
	server := mockServer(t, "capital_flow", map[string]interface{}{
		"symbol": "AAPL",
	})
	defer server.Close()

	qc := newTestQuoteClient(server.URL)
	data, err := qc.GetCapitalFlow("AAPL")
	if err != nil {
		t.Fatalf("GetCapitalFlow 失败: %v", err)
	}
	if data == nil {
		t.Fatal("GetCapitalFlow 返回 nil data")
	}
}

func TestGetCapitalDistribution(t *testing.T) {
	server := mockServer(t, "capital_distribution", map[string]interface{}{
		"symbol": "AAPL",
	})
	defer server.Close()

	qc := newTestQuoteClient(server.URL)
	data, err := qc.GetCapitalDistribution("AAPL")
	if err != nil {
		t.Fatalf("GetCapitalDistribution 失败: %v", err)
	}
	if data == nil {
		t.Fatal("GetCapitalDistribution 返回 nil data")
	}
}

// === 14.9 选股器和行情权限测试 ===

func TestMarketScanner(t *testing.T) {
	server := mockServer(t, "market_scanner", []map[string]interface{}{
		{"symbol": "AAPL"},
	})
	defer server.Close()

	qc := newTestQuoteClient(server.URL)
	params := map[string]interface{}{
		"market":   "US",
		"scanType": "TOP_GAINERS",
	}
	data, err := qc.MarketScanner(params)
	if err != nil {
		t.Fatalf("MarketScanner 失败: %v", err)
	}
	if data == nil {
		t.Fatal("MarketScanner 返回 nil data")
	}
}

func TestGrabQuotePermission(t *testing.T) {
	server := mockServer(t, "grab_quote_permission", []map[string]interface{}{
		{"name": "usStockQuote", "expireAt": 1700000000},
	})
	defer server.Close()

	qc := newTestQuoteClient(server.URL)
	data, err := qc.GrabQuotePermission()
	if err != nil {
		t.Fatalf("GrabQuotePermission 失败: %v", err)
	}
	if data == nil {
		t.Fatal("GrabQuotePermission 返回 nil data")
	}
}

// === 错误处理测试 ===

func TestQuoteClientApiError(t *testing.T) {
	server := mockErrorServer(t, 2100, "行情查询失败")
	defer server.Close()

	qc := newTestQuoteClient(server.URL)
	_, err := qc.GetMarketState("US")
	if err == nil {
		t.Fatal("期望返回错误，但返回了 nil")
	}
}
