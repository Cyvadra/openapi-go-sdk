// Package tigeropen 是老虎证券 OpenAPI Go SDK 的根包。
//
// 本 SDK 提供行情查询、交易下单、账户管理和实时推送等功能，
// 与 Python SDK 保持功能对等，遵循 Go 语言编码风格和最佳实践。
//
// 分层架构：
//   - 模型层（model）：Contract、Order、Position 等数据模型和枚举
//   - 配置层（config）：ClientConfig、ConfigParser
//   - 认证层（signer）：RSA 签名
//   - 传输层（client）：HttpClient、重试策略
//   - 业务层（quote/trade）：QuoteClient、TradeClient
//   - 推送层（push）：PushClient
package tigeropen

// Version 是 SDK 的版本号
const Version = "0.1.0"
