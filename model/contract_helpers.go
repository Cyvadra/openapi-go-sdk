package model

// StockContract 构造股票合约
func StockContract(symbol, currency string) Contract {
	return Contract{
		Symbol:   symbol,
		SecType:  string(SecTypeSTK),
		Currency: currency,
	}
}

// OptionContract 通过 identifier 构造期权合约
func OptionContract(identifier string) Contract {
	return Contract{
		SecType:    string(SecTypeOPT),
		Identifier: identifier,
	}
}

// OptionContractBySymbol 通过各字段构造期权合约
func OptionContractBySymbol(symbol, expiry string, strike float64, right, currency string) Contract {
	return Contract{
		Symbol:   symbol,
		SecType:  string(SecTypeOPT),
		Expiry:   expiry,
		Strike:   strike,
		Right:    right,
		Currency: currency,
	}
}

// FutureContract 构造期货合约
func FutureContract(symbol, currency, expiry string) Contract {
	return Contract{
		Symbol:   symbol,
		SecType:  string(SecTypeFUT),
		Currency: currency,
		Expiry:   expiry,
	}
}

// CashContract 构造外汇合约
func CashContract(symbol string) Contract {
	return Contract{
		Symbol:  symbol,
		SecType: string(SecTypeCASH),
	}
}

// FundContract 构造基金合约
func FundContract(symbol, currency string) Contract {
	return Contract{
		Symbol:   symbol,
		SecType:  string(SecTypeFUND),
		Currency: currency,
	}
}

// WarrantContract 构造窝轮合约
func WarrantContract(symbol, currency, expiry string, strike float64, right string) Contract {
	return Contract{
		Symbol:   symbol,
		SecType:  string(SecTypeWAR),
		Currency: currency,
		Expiry:   expiry,
		Strike:   strike,
		Right:    right,
	}
}
