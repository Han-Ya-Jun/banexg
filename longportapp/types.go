package longportapp

import (
	"github.com/banbox/banexg"
	"github.com/longportapp/openapi-go/config"
	"github.com/longportapp/openapi-go/quote"
	"github.com/longportapp/openapi-go/trade"
)

const (
	// Host keys
	HostQuote = "quote"
	HostTrade = "trade"
	
	// Quote API methods
	MethodQuoteGetSecurityList        = "QuoteGetSecurityList"
	MethodQuoteGetSecurityQuote       = "QuoteGetSecurityQuote"
	MethodQuoteGetSecurityDepth       = "QuoteGetSecurityDepth"
	MethodQuoteGetSecurityBrokers     = "QuoteGetSecurityBrokers"
	MethodQuoteGetSecurityTrades      = "QuoteGetSecurityTrades"
	MethodQuoteGetSecurityCandlesticks = "QuoteGetSecurityCandlesticks"
	
	// Trade API methods
	MethodTradeGetAccountBalance = "TradeGetAccountBalance"
	MethodTradeGetPositions      = "TradeGetPositions"
	MethodTradeGetOrders         = "TradeGetOrders"
	MethodTradeGetOrderDetail    = "TradeGetOrderDetail"
	MethodTradeSubmitOrder       = "TradeSubmitOrder"
	MethodTradeCancelOrder       = "TradeCancelOrder"
	MethodTradeReplaceOrder      = "TradeReplaceOrder"
)

// LongPortApp 长桥交易所实现
type LongPortApp struct {
	*banexg.Exchange
	
	// LongPort SDK clients
	quoteContext *quote.QuoteContext
	tradeContext *trade.TradeContext
	config       *config.Config
}

// LongPort 市场类型映射
const (
	MarketHK = "HK"  // 港股
	MarketUS = "US"  // 美股
	MarketCN = "CN"  // A股
)

// LongPort 订单类型
const (
	OrderTypeLO  = "LO"  // 限价单
	OrderTypeMO  = "MO"  // 市价单
	OrderTypeALO = "ALO" // 竞价限价单
	OrderTypeODD = "ODD" // 碎股单
	OrderTypeLIT = "LIT" // 增强限价单
	OrderTypeMIT = "MIT" // 市价转限价单
)

// LongPort 订单方向
const (
	OrderSideBuy  = "Buy"
	OrderSideSell = "Sell"
)

// LongPort 订单状态
const (
	OrderStatusSubmitted      = "Submitted"
	OrderStatusWaitingSubmit  = "WaitingSubmit"
	OrderStatusSubmitting     = "Submitting"
	OrderStatusSubmitFailed   = "SubmitFailed"
	OrderStatusFilled         = "Filled"
	OrderStatusPartialFilled  = "PartialFilled"
	OrderStatusCanceled       = "Canceled"
	OrderStatusCancelFailed   = "CancelFailed"
	OrderStatusRejected       = "Rejected"
)

// LongPort 时间周期
const (
	PeriodMin_1   = "1m"
	PeriodMin_5   = "5m"
	PeriodMin_15  = "15m"
	PeriodMin_30  = "30m"
	PeriodHour_1  = "1h"
	PeriodDay_1   = "1d"
	PeriodWeek_1  = "1w"
	PeriodMonth_1 = "1M"
)

// 默认关注的市场类型
var DefCareMarkets = []string{
	banexg.MarketSpot, // 将港股、美股、A股都映射为现货市场
}