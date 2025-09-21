# LongPort 长桥证券交易所适配器

本模块为banexg框架实现了LongPort长桥证券的API对接，支持港股、美股、A股的行情查询和交易功能。

## 🚀 功能特性

### 行情数据
- ✅ 获取单个/多个股票实时行情
- ✅ 获取订单簿深度数据
- ✅ 获取K线数据（支持多种时间周期）
- ✅ 获取历史交易数据

### 交易功能
- ✅ 获取账户余额
- ✅ 获取持仓信息
- ✅ 创建订单（限价单、市价单）
- ✅ 取消订单
- ✅ 查询订单详情
- ✅ 获取未完成订单

### 支持市场
- 🇭🇰 港股 (HK)
- 🇺🇸 美股 (US)  
- 🇨🇳 A股 (CN)

## 📦 安装配置

### 1. 获取API密钥

访问 [长桥开放平台](https://open.longportapp.com/) 申请API密钥：
- App Key
- App Secret  
- Access Token

### 2. 基本使用

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/banbox/banexg/bex"
)

func main() {
    // 配置API认证信息
    options := map[string]interface{}{
        "ApiKey":    "your_app_key",
        "ApiSecret": "your_app_secret", 
        "AccessToken":  "your_access_token", // AccessToken存储在password字段
    }

    // 创建交易所实例
    exg, err := bex.New("longportapp", options)
    if err != nil {
        log.Fatal("创建交易所失败:", err)
    }
    defer exg.Close()

    // 获取股票行情
    ticker, err := exg.FetchTicker("AAPL.US", nil)
    if err != nil {
        log.Fatal("获取行情失败:", err)
    }
    
    fmt.Printf("股票: %s, 最新价: %.2f\n", ticker.Symbol, ticker.Last)
}
```

### 3. 高级配置

```go
options := map[string]interface{}{
    "apiKey":      "your_app_key",
    "apiSecret":   "your_app_secret",
    "password":    "your_access_token",
    "marketType":  "spot",                    // 市场类型
    "careMarkets": []string{"spot"},          // 关注的市场
    "rateLimit":   100,                       // 请求间隔(ms)
    "timeout":     30000,                     // 超时时间(ms)
    "debugApi":    false,                     // 是否调试API
}
```

## 🔧 API 方法

### 行情数据

```go
// 获取单个股票行情
ticker, err := exg.FetchTicker("AAPL.US", nil)

// 获取多个股票行情
tickers, err := exg.FetchTickers([]string{"AAPL.US", "TSLA.US"}, nil)

// 获取订单簿
orderBook, err := exg.FetchOrderBook("AAPL.US", 10, nil)

// 获取K线数据
klines, err := exg.FetchOHLCV("AAPL.US", "1d", 0, 100, nil)
```

### 账户信息

```go
// 获取账户余额
balance, err := exg.FetchBalance(nil)

// 获取持仓信息
positions, err := exg.FetchPositions(nil, nil)
```

### 交易操作

```go
// 创建限价买单
order, err := exg.CreateOrder("AAPL.US", "limit", "buy", 100, 150.0, nil)

// 创建市价卖单
order, err := exg.CreateOrder("AAPL.US", "market", "sell", 100, 0, nil)

// 取消订单
canceledOrder, err := exg.CancelOrder(orderID, "AAPL.US", nil)

// 查询订单详情
orderDetail, err := exg.FetchOrder("AAPL.US", orderID, nil)

// 获取未完成订单
openOrders, err := exg.FetchOpenOrders("", 0, 10, nil)
```

## 📊 支持的时间周期

| 周期 | 说明 |
|------|------|
| 1m   | 1分钟 |
| 5m   | 5分钟 |
| 15m  | 15分钟 |
| 30m  | 30分钟 |
| 1h   | 1小时 |
| 1d   | 1天 |
| 1w   | 1周 |
| 1M   | 1月 |

## 🏷️ 股票代码格式

### 港股
- 腾讯控股: `700.HK`
- 阿里巴巴: `9988.HK`

### 美股  
- 苹果: `AAPL.US`
- 特斯拉: `TSLA.US`

### A股
- 平安银行: `000001.SZ`
- 贵州茅台: `600519.SH`

## ⚠️ 注意事项

1. **API限制**: 长桥API有请求频率限制，建议设置合适的rateLimit
2. **市场时间**: 股票市场有开盘和休市时间，非交易时间某些API可能返回空数据
3. **权限要求**: 交易相关功能需要相应的API权限
4. **测试环境**: 建议先在测试环境验证功能正常后再用于生产

## 🔗 相关链接

- [长桥开放平台](https://open.longportapp.com/)
- [API文档](https://open.longportapp.com/docs)
- [Go SDK](https://github.com/longportapp/openapi-go)
- [banexg框架](https://github.com/banbox/banexg)

## 📝 更新日志

### v1.0.0 (2024-01-XX)
- ✅ 实现基础行情数据获取
- ✅ 实现账户信息查询
- ✅ 实现基础交易功能
- ✅ 支持港股、美股、A股市场
- ✅ 集成到banexg框架

## 🤝 贡献

欢迎提交Issue和Pull Request来改进这个项目！

## 📄 许可证

本项目采用MIT许可证，详见LICENSE文件。