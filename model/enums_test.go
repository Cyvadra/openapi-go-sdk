package model

import (
	"testing"
)

// 测试 Market 枚举值
func TestMarketValues(t *testing.T) {
	tests := []struct {
		name string
		val  Market
		want string
	}{
		{"ALL", MarketAll, "ALL"},
		{"US", MarketUS, "US"},
		{"HK", MarketHK, "HK"},
		{"CN", MarketCN, "CN"},
		{"SG", MarketSG, "SG"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.val) != tt.want {
				t.Errorf("Market %s = %q, want %q", tt.name, tt.val, tt.want)
			}
		})
	}
}

// 测试 SecurityType 枚举值
func TestSecurityTypeValues(t *testing.T) {
	tests := []struct {
		name string
		val  SecurityType
		want string
	}{
		{"ALL", SecTypeAll, "ALL"},
		{"STK", SecTypeSTK, "STK"},
		{"OPT", SecTypeOPT, "OPT"},
		{"WAR", SecTypeWAR, "WAR"},
		{"IOPT", SecTypeIOPT, "IOPT"},
		{"FUT", SecTypeFUT, "FUT"},
		{"FOP", SecTypeFOP, "FOP"},
		{"CASH", SecTypeCASH, "CASH"},
		{"MLEG", SecTypeMLEG, "MLEG"},
		{"FUND", SecTypeFUND, "FUND"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.val) != tt.want {
				t.Errorf("SecurityType %s = %q, want %q", tt.name, tt.val, tt.want)
			}
		})
	}
}

// 测试 Currency 枚举值
func TestCurrencyValues(t *testing.T) {
	tests := []struct {
		name string
		val  Currency
		want string
	}{
		{"ALL", CurrencyAll, "ALL"},
		{"USD", CurrencyUSD, "USD"},
		{"HKD", CurrencyHKD, "HKD"},
		{"CNH", CurrencyCNH, "CNH"},
		{"SGD", CurrencySGD, "SGD"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.val) != tt.want {
				t.Errorf("Currency %s = %q, want %q", tt.name, tt.val, tt.want)
			}
		})
	}
}

// 测试 OrderType 枚举值
func TestOrderTypeValues(t *testing.T) {
	tests := []struct {
		name string
		val  OrderType
		want string
	}{
		{"MKT", OrderTypeMKT, "MKT"},
		{"LMT", OrderTypeLMT, "LMT"},
		{"STP", OrderTypeSTP, "STP"},
		{"STP_LMT", OrderTypeSTPLMT, "STP_LMT"},
		{"TRAIL", OrderTypeTRAIL, "TRAIL"},
		{"AM", OrderTypeAM, "AM"},
		{"AL", OrderTypeAL, "AL"},
		{"TWAP", OrderTypeTWAP, "TWAP"},
		{"VWAP", OrderTypeVWAP, "VWAP"},
		{"OCA", OrderTypeOCA, "OCA"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.val) != tt.want {
				t.Errorf("OrderType %s = %q, want %q", tt.name, tt.val, tt.want)
			}
		})
	}
}

// 测试 OrderStatus 枚举值
func TestOrderStatusValues(t *testing.T) {
	tests := []struct {
		name string
		val  OrderStatus
		want string
	}{
		{"PendingNew", OrderStatusPendingNew, "PendingNew"},
		{"Initial", OrderStatusInitial, "Initial"},
		{"Submitted", OrderStatusSubmitted, "Submitted"},
		{"PartiallyFilled", OrderStatusPartiallyFilled, "PartiallyFilled"},
		{"Filled", OrderStatusFilled, "Filled"},
		{"Cancelled", OrderStatusCancelled, "Cancelled"},
		{"PendingCancel", OrderStatusPendingCancel, "PendingCancel"},
		{"Inactive", OrderStatusInactive, "Inactive"},
		{"Invalid", OrderStatusInvalid, "Invalid"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.val) != tt.want {
				t.Errorf("OrderStatus %s = %q, want %q", tt.name, tt.val, tt.want)
			}
		})
	}
}

// 测试 BarPeriod 枚举值
func TestBarPeriodValues(t *testing.T) {
	tests := []struct {
		name string
		val  BarPeriod
		want string
	}{
		{"day", BarPeriodDay, "day"},
		{"week", BarPeriodWeek, "week"},
		{"month", BarPeriodMonth, "month"},
		{"year", BarPeriodYear, "year"},
		{"1min", BarPeriod1Min, "1min"},
		{"5min", BarPeriod5Min, "5min"},
		{"15min", BarPeriod15Min, "15min"},
		{"30min", BarPeriod30Min, "30min"},
		{"60min", BarPeriod60Min, "60min"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.val) != tt.want {
				t.Errorf("BarPeriod %s = %q, want %q", tt.name, tt.val, tt.want)
			}
		})
	}
}

// 测试 Language 枚举值
func TestLanguageValues(t *testing.T) {
	tests := []struct {
		name string
		val  Language
		want string
	}{
		{"zh_CN", LanguageZhCN, "zh_CN"},
		{"zh_TW", LanguageZhTW, "zh_TW"},
		{"en_US", LanguageEnUS, "en_US"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.val) != tt.want {
				t.Errorf("Language %s = %q, want %q", tt.name, tt.val, tt.want)
			}
		})
	}
}

// 测试 QuoteRight 枚举值
func TestQuoteRightValues(t *testing.T) {
	tests := []struct {
		name string
		val  QuoteRight
		want string
	}{
		{"br", QuoteRightBr, "br"},
		{"nr", QuoteRightNr, "nr"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.val) != tt.want {
				t.Errorf("QuoteRight %s = %q, want %q", tt.name, tt.val, tt.want)
			}
		})
	}
}

// 测试 License 枚举值
func TestLicenseValues(t *testing.T) {
	tests := []struct {
		name string
		val  License
		want string
	}{
		{"TBNZ", LicenseTBNZ, "TBNZ"},
		{"TBSG", LicenseTBSG, "TBSG"},
		{"TBHK", LicenseTBHK, "TBHK"},
		{"TBAU", LicenseTBAU, "TBAU"},
		{"TBUS", LicenseTBUS, "TBUS"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.val) != tt.want {
				t.Errorf("License %s = %q, want %q", tt.name, tt.val, tt.want)
			}
		})
	}
}

// 测试 TimeInForce 枚举值
func TestTimeInForceValues(t *testing.T) {
	tests := []struct {
		name string
		val  TimeInForce
		want string
	}{
		{"DAY", TimeInForceDAY, "DAY"},
		{"GTC", TimeInForceGTC, "GTC"},
		{"OPG", TimeInForceOPG, "OPG"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.val) != tt.want {
				t.Errorf("TimeInForce %s = %q, want %q", tt.name, tt.val, tt.want)
			}
		})
	}
}
