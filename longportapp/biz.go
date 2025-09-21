package longportapp

import (
	"context"
	"strconv"
	"time"

	"github.com/banbox/banexg"
	"github.com/banbox/banexg/errs"
	"github.com/longportapp/openapi-go/config"
	"github.com/longportapp/openapi-go/quote"
	"github.com/longportapp/openapi-go/trade"
	"github.com/shopspring/decimal"
)

func (e *LongPortApp) Init() *errs.Error {
	err := e.Exchange.Init()
	if err != nil {
		return err
	}

	if e.CareMarkets == nil || len(e.CareMarkets) == 0 {
		e.CareMarkets = DefCareMarkets
	}

	// 初始化LongPort配置
	err = e.initLongPortClients()
	if err != nil {
		return err
	}

	return nil
}

// 初始化LongPort客户端
func (e *LongPortApp) initLongPortClients() *errs.Error {
	// 从配置中获取认证信息
	accName := e.DefAccName
	if accName == "" {
		accName = "default"
	}

	acc, exists := e.Accounts[accName]
	if !exists || acc.Creds == nil {
		return errs.NewMsg(errs.CodeCredsRequired, "LongPort credentials required")
	}

	// 创建LongPort配置
	cfg := &config.Config{
		AppKey:      acc.Creds.ApiKey,
		AppSecret:   acc.Creds.Secret,
		AccessToken: acc.Creds.AccessToken, // 使用Password字段存储AccessToken
	}

	// 创建Quote客户端
	quoteCtx, err := quote.NewFromCfg(cfg)
	if err != nil {
		return errs.NewMsg(errs.CodeConnectFail, "failed to create quote context: %v", err)
	}

	// 创建Trade客户端
	tradeCtx, err := trade.NewFromCfg(cfg)
	if err != nil {
		quoteCtx.Close()
		return errs.NewMsg(errs.CodeConnectFail, "failed to create trade context: %v", err)
	}

	e.config = cfg
	e.quoteContext = quoteCtx
	e.tradeContext = tradeCtx

	return nil
}

// Close 关闭连接
func (e *LongPortApp) Close() *errs.Error {
	if e.quoteContext != nil {
		e.quoteContext.Close()
	}
	if e.tradeContext != nil {
		e.tradeContext.Close()
	}
	return e.Exchange.Close()
}

// FetchTicker 获取单个股票行情
func (e *LongPortApp) FetchTicker(symbol string, params map[string]interface{}) (*banexg.Ticker, *errs.Error) {
	if e.quoteContext == nil {
		return nil, errs.NewMsg(errs.CodeConnectFail, "quote context not initialized")
	}

	ctx := context.Background()
	securities := []string{symbol}

	quotes, err := e.quoteContext.Quote(ctx, securities)
	if err != nil {
		return nil, errs.NewMsg(errs.CodeRunTime, "failed to fetch quote: %v", err)
	}

	if len(quotes) == 0 {
		return nil, errs.NewMsg(errs.CodeRunTime, "no quote data for symbol: %s", symbol)
	}

	q := quotes[0]
	ticker := &banexg.Ticker{
		Symbol:        symbol,
		TimeStamp:     time.Now().UnixMilli(),
		Last:          q.LastDone.InexactFloat64(),
		Open:          q.Open.InexactFloat64(),
		High:          q.High.InexactFloat64(),
		Low:           q.Low.InexactFloat64(),
		Close:         q.LastDone.InexactFloat64(),
		BaseVolume:    float64(q.Volume),
		Change:        q.LastDone.Sub(*q.PrevClose).InexactFloat64(),
		Percentage:    q.LastDone.Sub(*q.PrevClose).Div(*q.PrevClose).Mul(decimal.NewFromInt(100)).InexactFloat64(),
		PreviousClose: q.PrevClose.InexactFloat64(),
	}

	return ticker, nil
}

// FetchTickers 获取多个股票行情
func (e *LongPortApp) FetchTickers(symbols []string, params map[string]interface{}) ([]*banexg.Ticker, *errs.Error) {
	if e.quoteContext == nil {
		return nil, errs.NewMsg(errs.CodeConnectFail, "quote context not initialized")
	}

	ctx := context.Background()
	quotes, err := e.quoteContext.Quote(ctx, symbols)
	if err != nil {
		return nil, errs.NewMsg(errs.CodeRunTime, "failed to fetch quotes: %v", err)
	}

	tickers := make([]*banexg.Ticker, 0, len(quotes))
	for i, q := range quotes {
		symbol := symbols[i]
		if i < len(symbols) {
			symbol = symbols[i]
		}

		ticker := &banexg.Ticker{
			Symbol:        symbol,
			TimeStamp:     time.Now().UnixMilli(),
			Last:          q.LastDone.InexactFloat64(),
			Open:          q.Open.InexactFloat64(),
			High:          q.High.InexactFloat64(),
			Low:           q.Low.InexactFloat64(),
			Close:         q.LastDone.InexactFloat64(),
			BaseVolume:    float64(q.Volume),
			Change:        q.LastDone.Sub(*q.PrevClose).InexactFloat64(),
			Percentage:    q.LastDone.Sub(*q.PrevClose).Div(*q.PrevClose).Mul(decimal.NewFromInt(100)).InexactFloat64(),
			PreviousClose: q.PrevClose.InexactFloat64(),
		}
		tickers = append(tickers, ticker)
	}

	return tickers, nil
}

// FetchOrderBook 获取订单簿
func (e *LongPortApp) FetchOrderBook(symbol string, limit int, params map[string]interface{}) (*banexg.OrderBook, *errs.Error) {
	if e.quoteContext == nil {
		return nil, errs.NewMsg(errs.CodeConnectFail, "quote context not initialized")
	}

	ctx := context.Background()
	depth, err := e.quoteContext.Depth(ctx, symbol)
	if err != nil {
		return nil, errs.NewMsg(errs.CodeRunTime, "failed to fetch depth: %v", err)
	}

	// 转换买盘数据
	bids := make([][2]float64, len(depth.Bid))
	for i, bid := range depth.Bid {
		bids[i] = [2]float64{
			bid.Price.InexactFloat64(),
			float64(bid.Volume),
		}
	}

	// 转换卖盘数据
	asks := make([][2]float64, len(depth.Ask))
	for i, ask := range depth.Ask {
		asks[i] = [2]float64{
			ask.Price.InexactFloat64(),
			float64(ask.Volume),
		}
	}

	orderBook := &banexg.OrderBook{
		Symbol:    symbol,
		TimeStamp: time.Now().UnixMilli(),
		Bids:      banexg.NewOdBookSide(true, limit, bids),
		Asks:      banexg.NewOdBookSide(false, limit, asks),
		Limit:     limit,
	}

	return orderBook, nil
}

// FetchOHLCV 获取K线数据
func (e *LongPortApp) FetchOHLCV(symbol, timeframe string, since int64, limit int, params map[string]interface{}) ([]*banexg.Kline, *errs.Error) {
	if e.quoteContext == nil {
		return nil, errs.NewMsg(errs.CodeConnectFail, "quote context not initialized")
	}

	// 转换时间周期
	period := e.convertTimeframe(timeframe)
	if period == quote.Period(0) {
		return nil, errs.NewMsg(errs.CodeParamInvalid, "unsupported timeframe: %s", timeframe)
	}

	ctx := context.Background()

	// 调用API - 使用正确的参数
	candlesticks, err := e.quoteContext.Candlesticks(ctx, symbol, period, int32(limit), quote.AdjustTypeNo)
	if err != nil {
		return nil, errs.NewMsg(errs.CodeRunTime, "failed to fetch candlesticks: %v", err)
	}

	klines := make([]*banexg.Kline, len(candlesticks))
	for i, candle := range candlesticks {
		klines[i] = &banexg.Kline{
			Time:   candle.Timestamp,
			Open:   candle.Open.InexactFloat64(),
			High:   candle.High.InexactFloat64(),
			Low:    candle.Low.InexactFloat64(),
			Close:  candle.Close.InexactFloat64(),
			Volume: float64(candle.Volume),
		}
	}

	return klines, nil
}

// convertTimeframe 转换时间周期格式
func (e *LongPortApp) convertTimeframe(timeframe string) quote.Period {
	switch timeframe {
	case "1m":
		return quote.PeriodOneMinute
	case "5m":
		return quote.PeriodFiveMinute
	case "15m":
		return quote.PeriodFifteenMinute
	case "30m":
		return quote.PeriodThirtyMinute
	case "1h":
		return quote.PeriodSixtyMinute
	case "1d":
		return quote.PeriodDay
	case "1w":
		return quote.PeriodWeek
	case "1M":
		return quote.PeriodMonth
	default:
		return quote.Period(0)
	}
}

// FetchBalance 获取账户余额
func (e *LongPortApp) FetchBalance(params map[string]interface{}) (*banexg.Balances, *errs.Error) {
	if e.tradeContext == nil {
		return nil, errs.NewMsg(errs.CodeConnectFail, "trade context not initialized")
	}

	ctx := context.Background()
	req := &trade.GetAccountBalance{}
	balances, err := e.tradeContext.AccountBalance(ctx, req)
	if err != nil {
		return nil, errs.NewMsg(errs.CodeRunTime, "failed to fetch balance: %v", err)
	}

	result := &banexg.Balances{
		TimeStamp: time.Now().UnixMilli(),
		Free:      make(map[string]float64),
		Used:      make(map[string]float64),
		Total:     make(map[string]float64),
		Assets:    make(map[string]*banexg.Asset),
	}

	for _, balance := range balances {
		currency := balance.Currency
		totalCash := balance.TotalCash.InexactFloat64()
		maxFinanceAmount := balance.MaxFinanceAmount.InexactFloat64()

		// 可用资金 = 总现金 - 最大融资金额
		free := totalCash - maxFinanceAmount
		if free < 0 {
			free = 0
		}

		result.Free[currency] = free
		result.Used[currency] = maxFinanceAmount
		result.Total[currency] = totalCash

		result.Assets[currency] = &banexg.Asset{
			Code:  currency,
			Free:  free,
			Used:  maxFinanceAmount,
			Total: totalCash,
		}
	}

	return result.Init(), nil
}

// FetchPositions 获取持仓信息
func (e *LongPortApp) FetchPositions(symbols []string, params map[string]interface{}) ([]*banexg.Position, *errs.Error) {
	if e.tradeContext == nil {
		return nil, errs.NewMsg(errs.CodeConnectFail, "trade context not initialized")
	}

	ctx := context.Background()
	positionChannels, err := e.tradeContext.StockPositions(ctx, symbols)
	if err != nil {
		return nil, errs.NewMsg(errs.CodeRunTime, "failed to fetch positions: %v", err)
	}

	result := make([]*banexg.Position, 0)
	for _, posChannel := range positionChannels {
		for _, pos := range posChannel.Positions {
			// 过滤指定的symbols
			if len(symbols) > 0 {
				found := false
				for _, symbol := range symbols {
					if pos.Symbol == symbol {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}

			quantity, _ := strconv.ParseFloat(pos.Quantity, 64)

			position := &banexg.Position{
				Symbol:     pos.Symbol,
				TimeStamp:  time.Now().UnixMilli(),
				Contracts:  quantity,
				EntryPrice: pos.CostPrice.InexactFloat64(),
				Side:       e.convertPositionSide(int64(quantity)),
			}

			result = append(result, position)
		}
	}

	return result, nil
}

// convertPositionSide 转换持仓方向
func (e *LongPortApp) convertPositionSide(quantity int64) string {
	if quantity > 0 {
		return banexg.PosSideLong
	} else if quantity < 0 {
		return banexg.PosSideShort
	}
	return banexg.PosSideBoth
}

// CreateOrder 创建订单
func (e *LongPortApp) CreateOrder(symbol, odType, side string, amount, price float64, params map[string]interface{}) (*banexg.Order, *errs.Error) {
	if e.tradeContext == nil {
		return nil, errs.NewMsg(errs.CodeConnectFail, "trade context not initialized")
	}

	// 转换订单类型
	orderType := e.convertOrderType(odType)
	if orderType == "" {
		return nil, errs.NewMsg(errs.CodeParamInvalid, "unsupported order type: %s", odType)
	}

	// 转换订单方向
	orderSide := e.convertOrderSide(side)
	if orderSide == "" {
		return nil, errs.NewMsg(errs.CodeParamInvalid, "unsupported order side: %s", side)
	}

	ctx := context.Background()

	// 构建订单请求
	order := &trade.SubmitOrder{
		Symbol:            symbol,
		OrderType:         trade.OrderType(orderType),
		Side:              trade.OrderSide(orderSide),
		SubmittedQuantity: uint64(amount),
		TimeInForce:       trade.TimeTypeDay, // 默认当日有效
	}

	// 如果是限价单，设置价格
	if odType == banexg.OdTypeLimit {
		order.SubmittedPrice = decimal.NewFromFloat(price)
	}

	// 提交订单
	orderID, err := e.tradeContext.SubmitOrder(ctx, order)
	if err != nil {
		return nil, errs.NewMsg(errs.CodeRunTime, "failed to submit order: %v", err)
	}

	// 返回订单信息
	result := &banexg.Order{
		ID:        orderID,
		Symbol:    symbol,
		Type:      odType,
		Side:      side,
		Amount:    amount,
		Price:     price,
		Status:    banexg.OdStatusOpen,
		Timestamp: time.Now().UnixMilli(),
	}

	return result, nil
}

// convertOrderType 转换订单类型
func (e *LongPortApp) convertOrderType(odType string) string {
	switch odType {
	case banexg.OdTypeLimit:
		return OrderTypeLO
	case banexg.OdTypeMarket:
		return OrderTypeMO
	default:
		return ""
	}
}

// convertOrderSide 转换订单方向
func (e *LongPortApp) convertOrderSide(side string) string {
	switch side {
	case banexg.OdSideBuy:
		return OrderSideBuy
	case banexg.OdSideSell:
		return OrderSideSell
	default:
		return ""
	}
}

// CancelOrder 取消订单
func (e *LongPortApp) CancelOrder(id string, symbol string, params map[string]interface{}) (*banexg.Order, *errs.Error) {
	if e.tradeContext == nil {
		return nil, errs.NewMsg(errs.CodeConnectFail, "trade context not initialized")
	}

	ctx := context.Background()
	err := e.tradeContext.CancelOrder(ctx, id)
	if err != nil {
		return nil, errs.NewMsg(errs.CodeRunTime, "failed to cancel order: %v", err)
	}

	// 返回取消的订单信息
	result := &banexg.Order{
		ID:        id,
		Symbol:    symbol,
		Status:    banexg.OdStatusCanceled,
		Timestamp: time.Now().UnixMilli(),
	}

	return result, nil
}

// FetchOrder 获取订单详情
func (e *LongPortApp) FetchOrder(symbol, orderId string, params map[string]interface{}) (*banexg.Order, *errs.Error) {
	if e.tradeContext == nil {
		return nil, errs.NewMsg(errs.CodeConnectFail, "trade context not initialized")
	}

	ctx := context.Background()
	orderDetail, err := e.tradeContext.OrderDetail(ctx, orderId)
	if err != nil {
		return nil, errs.NewMsg(errs.CodeRunTime, "failed to fetch order: %v", err)
	}

	return e.convertOrderDetailToOrder(&orderDetail), nil
}

// FetchOpenOrders 获取未完成订单
func (e *LongPortApp) FetchOpenOrders(symbol string, since int64, limit int, params map[string]interface{}) ([]*banexg.Order, *errs.Error) {
	if e.tradeContext == nil {
		return nil, errs.NewMsg(errs.CodeConnectFail, "trade context not initialized")
	}

	ctx := context.Background()
	req := &trade.GetTodayOrders{}

	if symbol != "" {
		req.Symbol = symbol
	}

	orders, err := e.tradeContext.TodayOrders(ctx, req)
	if err != nil {
		return nil, errs.NewMsg(errs.CodeRunTime, "failed to fetch open orders: %v", err)
	}

	result := make([]*banexg.Order, 0, len(orders))
	for _, order := range orders {
		// 只返回未完成的订单
		if e.isOrderOpen(order.Status) {
			result = append(result, e.convertToOrder(order))
		}
	}

	return result, nil
}

// isOrderOpen 判断订单是否未完成
func (e *LongPortApp) isOrderOpen(status trade.OrderStatus) bool {
	switch status {
	case trade.OrderNewStatus, trade.OrderWaitToNew,
		trade.OrderPartialFilledStatus:
		return true
	default:
		return false
	}
}

// convertToOrder 转换为标准订单格式
func (e *LongPortApp) convertToOrder(order *trade.Order) *banexg.Order {
	quantity, _ := strconv.ParseFloat(order.Quantity, 64)
	executedQuantity, _ := strconv.ParseFloat(order.ExecutedQuantity, 64)

	result := &banexg.Order{
		ID:        order.OrderId,
		Symbol:    order.Symbol,
		Type:      e.convertFromOrderType(string(order.OrderType)),
		Side:      e.convertFromOrderSide(string(order.Side)),
		Amount:    quantity,
		Filled:    executedQuantity,
		Remaining: quantity - executedQuantity,
		Status:    e.convertOrderStatus(string(order.Status)),
		Timestamp: e.parseTimeToMilli(order.SubmittedAt),
	}

	if order.Price != nil {
		result.Price = order.Price.InexactFloat64()
	}

	if order.ExecutedPrice != nil {
		result.Cost = order.ExecutedPrice.Mul(decimal.NewFromFloat(executedQuantity)).InexactFloat64()
	}

	return result
}

// convertFromOrderType 从LongPort订单类型转换为标准类型
func (e *LongPortApp) convertFromOrderType(orderType string) string {
	switch orderType {
	case OrderTypeLO:
		return banexg.OdTypeLimit
	case OrderTypeMO:
		return banexg.OdTypeMarket
	default:
		return orderType
	}
}

// convertFromOrderSide 从LongPort订单方向转换为标准方向
func (e *LongPortApp) convertFromOrderSide(side string) string {
	switch side {
	case OrderSideBuy:
		return banexg.OdSideBuy
	case OrderSideSell:
		return banexg.OdSideSell
	default:
		return side
	}
}

// convertOrderStatus 转换订单状态
func (e *LongPortApp) convertOrderStatus(status string) string {
	switch status {
	case string(trade.OrderNewStatus), string(trade.OrderWaitToNew):
		return banexg.OdStatusOpen
	case string(trade.OrderPartialFilledStatus):
		return banexg.OdStatusPartFilled
	case string(trade.OrderFilledStatus):
		return banexg.OdStatusFilled
	case string(trade.OrderCanceledStatus):
		return banexg.OdStatusCanceled
	case string(trade.OrderRejectedStatus):
		return banexg.OdStatusRejected
	default:
		return status
	}
}

// parseTimeToMilli 解析时间字符串为毫秒时间戳
func (e *LongPortApp) parseTimeToMilli(timeStr string) int64 {
	if timeStr == "" {
		return time.Now().UnixMilli()
	}

	// 尝试解析常见的时间格式
	layouts := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, timeStr); err == nil {
			return t.UnixMilli()
		}
	}

	// 如果解析失败，返回当前时间
	return time.Now().UnixMilli()
}

// convertOrderDetailToOrder 转换 OrderDetail 为标准订单格式
func (e *LongPortApp) convertOrderDetailToOrder(orderDetail *trade.OrderDetail) *banexg.Order {
	quantity := float64(orderDetail.Quantity)
	executedQuantity := float64(orderDetail.ExecutedQuantity)

	result := &banexg.Order{
		ID:        orderDetail.OrderId,
		Symbol:    orderDetail.Symbol,
		Type:      e.convertFromOrderType(string(orderDetail.OrderType)),
		Side:      e.convertFromOrderSide(string(orderDetail.Side)),
		Amount:    quantity,
		Filled:    executedQuantity,
		Remaining: quantity - executedQuantity,
		Status:    e.convertOrderStatus(string(orderDetail.Status)),
		Timestamp: e.parseTimeToMilli(orderDetail.SubmittedAt),
	}

	if orderDetail.Price != nil {
		result.Price = orderDetail.Price.InexactFloat64()
	}

	if orderDetail.ExecutedPrice != nil {
		result.Cost = orderDetail.ExecutedPrice.Mul(decimal.NewFromFloat(executedQuantity)).InexactFloat64()
	}

	return result
}
