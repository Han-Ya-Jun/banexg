package main

import (
	"fmt"
	"log"

	"github.com/banbox/banexg"
	"github.com/banbox/banexg/bex"
)

func main() {
	// 配置LongPort API凭证
	options := map[string]interface{}{
		"ApiKey":   "your_app_key_here",     // 替换为你的App Key
		"Secret":   "your_app_secret_here",  // 替换为你的App Secret
		"Password": "your_access_token_here", // 替换为你的Access Token
	}

	// 创建LongPort交易所实例
	exchange, err := bex.New("longportapp", options)
	if err != nil {
		log.Fatal("创建交易所失败:", err)
	}
	defer exchange.Close()

	fmt.Println("=== LongPort 长桥证券 API 示例 ===")

	// 1. 获取腾讯控股的实时报价
	fmt.Println("\n1. 获取腾讯控股(700.HK)实时报价:")
	ticker, err := exchange.FetchTicker("700.HK", nil)
	if err != nil {
		log.Printf("获取报价失败: %v", err)
	} else {
		fmt.Printf("股票代码: %s\n", ticker.Symbol)
		fmt.Printf("最新价格: %.2f\n", ticker.Last)
		fmt.Printf("开盘价: %.2f\n", ticker.Open)
		fmt.Printf("最高价: %.2f\n", ticker.High)
		fmt.Printf("最低价: %.2f\n", ticker.Low)
		fmt.Printf("涨跌额: %.2f\n", ticker.Change)
		fmt.Printf("涨跌幅: %.2f%%\n", ticker.Percentage)
		fmt.Printf("成交量: %.0f\n", ticker.Volume)
	}

	// 2. 获取多个股票的报价
	fmt.Println("\n2. 获取多个股票报价:")
	symbols := []string{"700.HK", "AAPL.US", "000001.SZ"}
	tickers, err := exchange.FetchTickers(symbols, nil)
	if err != nil {
		log.Printf("获取多个报价失败: %v", err)
	} else {
		for _, t := range tickers {
			fmt.Printf("%s: %.2f (%.2f%%)\n", t.Symbol, t.Last, t.Percentage)
		}
	}

	// 3. 获取订单簿深度
	fmt.Println("\n3. 获取腾讯控股订单簿深度:")
	orderBook, err := exchange.FetchOrderBook("700.HK", 5, nil)
	if err != nil {
		log.Printf("获取订单簿失败: %v", err)
	} else {
		fmt.Println("买盘 (Bids):")
		for i, bid := range orderBook.Bids.Price {
			if i >= 5 {
				break
			}
			fmt.Printf("  价格: %.2f, 数量: %.0f\n", bid, orderBook.Bids.Volume[i])
		}
		fmt.Println("卖盘 (Asks):")
		for i, ask := range orderBook.Asks.Price {
			if i >= 5 {
				break
			}
			fmt.Printf("  价格: %.2f, 数量: %.0f\n", ask, orderBook.Asks.Volume[i])
		}
	}

	// 4. 获取K线数据
	fmt.Println("\n4. 获取腾讯控股日K线数据(最近5天):")
	klines, err := exchange.FetchOHLCV("700.HK", "1d", 0, 5, nil)
	if err != nil {
		log.Printf("获取K线失败: %v", err)
	} else {
		fmt.Println("时间\t\t开盘\t最高\t最低\t收盘\t成交量")
		for _, k := range klines {
			fmt.Printf("%d\t%.2f\t%.2f\t%.2f\t%.2f\t%.0f\n",
				k.Time, k.Open, k.High, k.Low, k.Close, k.Volume)
		}
	}

	// 5. 获取账户余额
	fmt.Println("\n5. 获取账户余额:")
	balance, err := exchange.FetchBalance(nil)
	if err != nil {
		log.Printf("获取余额失败: %v", err)
	} else {
		fmt.Println("货币\t\t可用\t\t冻结\t\t总计")
		for currency, total := range balance.Total {
			free := balance.Free[currency]
			used := balance.Used[currency]
			fmt.Printf("%s\t\t%.2f\t\t%.2f\t\t%.2f\n", currency, free, used, total)
		}
	}

	// 6. 获取持仓信息
	fmt.Println("\n6. 获取持仓信息:")
	positions, err := exchange.FetchPositions(nil, nil)
	if err != nil {
		log.Printf("获取持仓失败: %v", err)
	} else {
		if len(positions) == 0 {
			fmt.Println("暂无持仓")
		} else {
			fmt.Println("股票代码\t\t数量\t\t成本价\t\t市价\t\t盈亏")
			for _, pos := range positions {
				fmt.Printf("%s\t\t%.0f\t\t%.2f\t\t%.2f\t\t%.2f\n",
					pos.Symbol, pos.Contracts, pos.EntryPrice, pos.MarkPrice, pos.UnrealizedPnl)
			}
		}
	}

	// 7. 创建测试订单 (注意：这会创建真实订单，请谨慎使用)
	fmt.Println("\n7. 创建测试订单 (已注释，取消注释前请确认):")
	fmt.Println("// 创建腾讯控股限价买入订单")
	fmt.Println("// order, err := exchange.CreateOrder(\"700.HK\", banexg.OdTypeLimit, banexg.OdSideBuy, 100, 400.0, nil)")
	
	/*
	// 取消注释以下代码来创建真实订单 (请谨慎操作)
	order, err := exchange.CreateOrder("700.HK", banexg.OdTypeLimit, banexg.OdSideBuy, 100, 400.0, nil)
	if err != nil {
		log.Printf("创建订单失败: %v", err)
	} else {
		fmt.Printf("订单创建成功: %s\n", order.ID)
		
		// 查询订单状态
		orderDetail, err := exchange.FetchOrder("700.HK", order.ID, nil)
		if err != nil {
			log.Printf("查询订单失败: %v", err)
		} else {
			fmt.Printf("订单状态: %s\n", orderDetail.Status)
		}
		
		// 取消订单
		canceledOrder, err := exchange.CancelOrder(order.ID, "700.HK", nil)
		if err != nil {
			log.Printf("取消订单失败: %v", err)
		} else {
			fmt.Printf("订单已取消: %s\n", canceledOrder.ID)
		}
	}
	*/

	// 8. 获取未完成订单
	fmt.Println("\n8. 获取未完成订单:")
	openOrders, err := exchange.FetchOpenOrders("", 0, 0, nil)
	if err != nil {
		log.Printf("获取未完成订单失败: %v", err)
	} else {
		if len(openOrders) == 0 {
			fmt.Println("暂无未完成订单")
		} else {
			fmt.Println("订单ID\t\t股票代码\t类型\t方向\t数量\t价格\t状态")
			for _, order := range openOrders {
				fmt.Printf("%s\t%s\t%s\t%s\t%.0f\t%.2f\t%s\n",
					order.ID, order.Symbol, order.Type, order.Side,
					order.Amount, order.Price, order.Status)
			}
		}
	}

	fmt.Println("\n=== 示例完成 ===")
	fmt.Println("注意：")
	fmt.Println("1. 请替换示例中的API凭证为你的真实凭证")
	fmt.Println("2. 股票代码格式：港股(700.HK)、美股(AAPL.US)、A股(000001.SZ)")
	fmt.Println("3. 交易功能需要在交易时间内使用")
	fmt.Println("4. 创建订单前请确认参数，避免误操作")
}