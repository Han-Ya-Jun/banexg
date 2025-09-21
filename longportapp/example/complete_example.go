package main

import (
	"fmt"
	"log"
	"time"

	"github.com/banbox/banexg/bex"
	"github.com/banbox/banexg/longportapp"
)

func main() {
	// 配置LongPort API认证信息
	options := map[string]interface{}{
		"apiKey":      "your_app_key",
		"apiSecret":   "your_app_secret", 
		"password":    "your_access_token", // AccessToken存储在password字段
		"marketType":  "spot",
		"careMarkets": []string{"spot"},
	}

	// 创建LongPort交易所实例
	exg, err := bex.New("longportapp", options)
	if err != nil {
		log.Fatal("创建交易所失败:", err)
	}
	defer exg.Close()

	// 转换为LongPortApp类型以使用特定功能
	longport := exg.(*longportapp.LongPortApp)

	// 示例1: 获取单个股票行情
	fmt.Println("=== 获取单个股票行情 ===")
	ticker, err := longport.FetchTicker("AAPL.US", nil)
	if err != nil {
		log.Printf("获取行情失败: %v", err)
	} else {
		fmt.Printf("股票: %s, 最新价: %.2f, 涨跌: %.2f (%.2f%%)\n", 
			ticker.Symbol, ticker.Last, ticker.Change, ticker.Percentage)
	}

	// 示例2: 获取多个股票行情
	fmt.Println("\n=== 获取多个股票行情 ===")
	symbols := []string{"AAPL.US", "TSLA.US", "700.HK"}
	tickers, err := longport.FetchTickers(symbols, nil)
	if err != nil {
		log.Printf("获取多个行情失败: %v", err)
	} else {
		for _, ticker := range tickers {
			fmt.Printf("股票: %s, 最新价: %.2f, 涨跌: %.2f (%.2f%%)\n", 
				ticker.Symbol, ticker.Last, ticker.Change, ticker.Percentage)
		}
	}

	// 示例3: 获取订单簿
	fmt.Println("\n=== 获取订单簿 ===")
	orderBook, err := longport.FetchOrderBook("AAPL.US", 10, nil)
	if err != nil {
		log.Printf("获取订单簿失败: %v", err)
	} else {
		fmt.Printf("订单簿 - 买盘数量: %d, 卖盘数量: %d\n", 
			len(orderBook.Bids.Items), len(orderBook.Asks.Items))
		if len(orderBook.Bids.Items) > 0 {
			fmt.Printf("最佳买价: %.2f, 数量: %.0f\n", 
				orderBook.Bids.Items[0].Price, orderBook.Bids.Items[0].Size)
		}
		if len(orderBook.Asks.Items) > 0 {
			fmt.Printf("最佳卖价: %.2f, 数量: %.0f\n", 
				orderBook.Asks.Items[0].Price, orderBook.Asks.Items[0].Size)
		}
	}

	// 示例4: 获取K线数据
	fmt.Println("\n=== 获取K线数据 ===")
	klines, err := longport.FetchOHLCV("AAPL.US", "1d", 0, 5, nil)
	if err != nil {
		log.Printf("获取K线失败: %v", err)
	} else {
		fmt.Printf("获取到 %d 根K线\n", len(klines))
		for i, kline := range klines {
			fmt.Printf("K线%d: 时间=%s, 开盘=%.2f, 最高=%.2f, 最低=%.2f, 收盘=%.2f, 成交量=%.0f\n",
				i+1, time.Unix(kline.Time/1000, 0).Format("2006-01-02"), 
				kline.Open, kline.High, kline.Low, kline.Close, kline.Volume)
		}
	}

	// 示例5: 获取账户余额
	fmt.Println("\n=== 获取账户余额 ===")
	balance, err := longport.FetchBalance(nil)
	if err != nil {
		log.Printf("获取余额失败: %v", err)
	} else {
		fmt.Printf("账户余额信息:\n")
		for currency, asset := range balance.Assets {
			if asset.Total > 0 {
				fmt.Printf("  %s: 总计=%.2f, 可用=%.2f, 冻结=%.2f\n", 
					currency, asset.Total, asset.Free, asset.Used)
			}
		}
	}

	// 示例6: 获取持仓信息
	fmt.Println("\n=== 获取持仓信息 ===")
	positions, err := longport.FetchPositions(nil, nil)
	if err != nil {
		log.Printf("获取持仓失败: %v", err)
	} else {
		fmt.Printf("持仓数量: %d\n", len(positions))
		for _, pos := range positions {
			fmt.Printf("  股票: %s, 数量: %.0f, 成本价: %.2f, 市价: %.2f, 浮盈: %.2f\n",
				pos.Symbol, pos.Contracts, pos.EntryPrice, pos.MarkPrice, pos.UnrealizedPnl)
		}
	}

	// 示例7: 创建限价买单 (注意: 这会实际下单，请谨慎使用)
	fmt.Println("\n=== 创建订单示例 (已注释，取消注释前请确认) ===")
	/*
	order, err := longport.CreateOrder("AAPL.US", "limit", "buy", 100, 150.0, nil)
	if err != nil {
		log.Printf("创建订单失败: %v", err)
	} else {
		fmt.Printf("订单创建成功: ID=%s, 股票=%s, 类型=%s, 方向=%s, 数量=%.0f, 价格=%.2f\n",
			order.ID, order.Symbol, order.Type, order.Side, order.Amount, order.Price)
		
		// 获取订单详情
		orderDetail, err := longport.FetchOrder(order.Symbol, order.ID, nil)
		if err != nil {
			log.Printf("获取订单详情失败: %v", err)
		} else {
			fmt.Printf("订单详情: 状态=%s, 已成交=%.0f, 剩余=%.0f\n",
				orderDetail.Status, orderDetail.Filled, orderDetail.Remaining)
		}
		
		// 取消订单
		canceledOrder, err := longport.CancelOrder(order.ID, order.Symbol, nil)
		if err != nil {
			log.Printf("取消订单失败: %v", err)
		} else {
			fmt.Printf("订单已取消: ID=%s, 状态=%s\n", canceledOrder.ID, canceledOrder.Status)
		}
	}
	*/

	// 示例8: 获取未完成订单
	fmt.Println("\n=== 获取未完成订单 ===")
	openOrders, err := longport.FetchOpenOrders("", 0, 10, nil)
	if err != nil {
		log.Printf("获取未完成订单失败: %v", err)
	} else {
		fmt.Printf("未完成订单数量: %d\n", len(openOrders))
		for _, order := range openOrders {
			fmt.Printf("  订单: ID=%s, 股票=%s, 类型=%s, 方向=%s, 数量=%.0f, 价格=%.2f, 状态=%s\n",
				order.ID, order.Symbol, order.Type, order.Side, order.Amount, order.Price, order.Status)
		}
	}

	fmt.Println("\n=== 示例完成 ===")
}