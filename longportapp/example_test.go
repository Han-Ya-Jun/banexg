package longportapp

import (
	"fmt"
	"log"
	"testing"

	"github.com/banbox/banexg"
)

// 示例：如何使用LongPort交易所适配器
func ExampleLongPortApp() {
	// 创建交易所实例
	options := map[string]interface{}{
		"ApiKey":      "xxxx",
		"ApiSecret":   "xxxxx",
		"AccessToken": "xxxxx",
	}

	exchange, err := New(options)
	if err != nil {
		log.Fatal(err)
	}
	defer exchange.Close()

	// 获取腾讯控股的实时报价
	ticker, err := exchange.FetchTicker("700.HK", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("腾讯控股: 价格=%f, 涨跌=%f%%\n", ticker.Last, ticker.Percentage)

	// 获取多个股票的报价
	symbols := []string{"700.HK", "AAPL.US", "000001.SZ"}
	tickers, err := exchange.FetchTickers(symbols, nil)
	if err != nil {
		log.Fatal(err)
	}
	for _, t := range tickers {
		fmt.Printf("%s: %f\n", t.Symbol, t.Last)
	}

	// 获取订单簿
	orderBook, err := exchange.FetchOrderBook("700.HK", 5, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("订单簿深度: 买盘%d档, 卖盘%d档\n", len(orderBook.Bids.Price), len(orderBook.Asks.Price))

	// 获取K线数据
	klines, err := exchange.FetchOHLCV("700.HK", "1d", 0, 10, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("获取到%d根K线\n", len(klines))

	// 获取账户余额
	balance, err := exchange.FetchBalance(nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("账户余额: %+v\n", balance.Total)

	// 获取持仓
	positions, err := exchange.FetchPositions(nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("持仓数量: %d\n", len(positions))

	// 创建限价买入订单
	order, err := exchange.CreateOrder(
		"700.HK",           // 腾讯控股
		banexg.OdTypeLimit, // 限价单
		banexg.OdSideBuy,   // 买入
		100,                // 100股
		400.0,              // 400港币
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("订单创建成功: %s\n", order.ID)

	// 查询订单状态
	orderDetail, err := exchange.FetchOrder("700.HK", order.ID, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("订单状态: %s\n", orderDetail.Status)

	// 获取未完成订单
	openOrders, err := exchange.FetchOpenOrders("", 0, 0, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("未完成订单数量: %d\n", len(openOrders))

	// 取消订单
	canceledOrder, err := exchange.CancelOrder(order.ID, "700.HK", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("订单已取消: %s\n", canceledOrder.ID)
}

func TestLongPortAppBasic(t *testing.T) {
	ExampleLongPortApp()
}

//// 测试基本功能
//func TestLongPortAppBasic(t *testing.T) {
//	// 注意：这个测试需要真实的API凭证才能运行
//	t.Skip("需要真实的API凭证")
//
//	options := map[string]interface{}{
//		"ApiKey":   "test_app_key",
//		"Secret":   "test_app_secret",
//		"Password": "test_access_token",
//	}
//
//	exchange, err := New(options)
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer exchange.Close()
//
//	// 测试获取市场信息
//	markets, err := exchange.LoadMarkets(false, nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//	if len(markets) == 0 {
//		t.Error("应该有市场数据")
//	}
//
//	// 测试获取报价（使用模拟数据）
//	// 实际测试时需要替换为真实的股票代码
//	ticker, err := exchange.FetchTicker("700.HK", nil)
//	if err != nil {
//		t.Log("获取报价失败（可能是网络或认证问题）:", err)
//	} else {
//		if ticker.Symbol != "700.HK" {
//			t.Error("股票代码不匹配")
//		}
//		if ticker.Last <= 0 {
//			t.Error("价格应该大于0")
//		}
//	}
//}
//
//// 基准测试
//func BenchmarkFetchTicker(b *testing.B) {
//	b.Skip("需要真实的API凭证")
//
//	options := map[string]interface{}{
//		"ApiKey":   "test_app_key",
//		"Secret":   "test_app_secret",
//		"Password": "test_access_token",
//	}
//
//	exchange, err := New(options)
//	if err != nil {
//		b.Fatal(err)
//	}
//	defer exchange.Close()
//
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		_, err := exchange.FetchTicker("700.HK", nil)
//		if err != nil {
//			b.Fatal(err)
//		}
//	}
//}
