package longportapp

import (
	"github.com/banbox/banexg"
	"github.com/banbox/banexg/errs"
)

func New(Options map[string]interface{}) (*LongPortApp, *errs.Error) {
	exg := &LongPortApp{
		Exchange: &banexg.Exchange{
			ExgInfo: &banexg.ExgInfo{
				ID:        "longportapp",
				Name:      "LongPort",
				Countries: []string{"HK", "US", "CN"},
				NoHoliday: false, // 股票市场有休市时间
				FullDay:   false, // 股票市场不是24小时交易
			},
			RateLimit: 100, // 100ms间隔
			Options:   Options,
			Hosts: &banexg.ExgHosts{
				Test: map[string]string{
					HostQuote: "https://openapi.longportapp.com",
					HostTrade: "https://openapi.longportapp.com",
				},
				Prod: map[string]string{
					HostQuote: "https://openapi.longportapp.com",
					HostTrade: "https://openapi.longportapp.com",
				},
				Www: "https://longportapp.com",
				Doc: []string{
					"https://open.longportapp.com/docs",
				},
				Fees: "https://longportapp.com/fees",
			},
			Fees: &banexg.ExgFee{
				Main: &banexg.TradeFee{
					FeeSide:    "quote",
					TierBased:  false,
					Percentage: true,
					Taker:      0.003, // 0.3% 默认费率
					Maker:      0.003, // 0.3% 默认费率
				},
			},
			Apis: map[string]*banexg.Entry{
				// Quote APIs
				MethodQuoteGetSecurityList:    {Path: "v1/quote/security/list", Host: HostQuote, Method: "GET", Cost: 1},
				MethodQuoteGetSecurityQuote:   {Path: "v1/quote/security/quote", Host: HostQuote, Method: "GET", Cost: 1},
				MethodQuoteGetSecurityDepth:   {Path: "v1/quote/security/depth", Host: HostQuote, Method: "GET", Cost: 1},
				MethodQuoteGetSecurityBrokers: {Path: "v1/quote/security/brokers", Host: HostQuote, Method: "GET", Cost: 1},
				MethodQuoteGetSecurityTrades:  {Path: "v1/quote/security/trades", Host: HostQuote, Method: "GET", Cost: 1},
				MethodQuoteGetSecurityCandlesticks: {Path: "v1/quote/security/candlesticks", Host: HostQuote, Method: "GET", Cost: 1},
				
				// Trade APIs
				MethodTradeGetAccountBalance: {Path: "v1/trade/account/balance", Host: HostTrade, Method: "GET", Cost: 1},
				MethodTradeGetPositions:      {Path: "v1/trade/position/list", Host: HostTrade, Method: "GET", Cost: 1},
				MethodTradeGetOrders:         {Path: "v1/trade/order/list", Host: HostTrade, Method: "GET", Cost: 1},
				MethodTradeGetOrderDetail:    {Path: "v1/trade/order/detail", Host: HostTrade, Method: "GET", Cost: 1},
				MethodTradeSubmitOrder:       {Path: "v1/trade/order/submit", Host: HostTrade, Method: "POST", Cost: 1},
				MethodTradeCancelOrder:       {Path: "v1/trade/order/cancel", Host: HostTrade, Method: "DELETE", Cost: 1},
				MethodTradeReplaceOrder:      {Path: "v1/trade/order/replace", Host: HostTrade, Method: "PUT", Cost: 1},
			},
			Has: map[string]map[string]int{
				"": {
					banexg.ApiFetchTicker:           banexg.HasOk,
					banexg.ApiFetchTickers:          banexg.HasOk,
					banexg.ApiFetchOHLCV:            banexg.HasOk,
					banexg.ApiFetchOrderBook:        banexg.HasOk,
					banexg.ApiFetchOrder:            banexg.HasOk,
					banexg.ApiFetchOrders:           banexg.HasOk,
					banexg.ApiFetchBalance:          banexg.HasOk,
					banexg.ApiFetchPositions:        banexg.HasOk,
					banexg.ApiFetchOpenOrders:       banexg.HasOk,
					banexg.ApiCreateOrder:           banexg.HasOk,
					banexg.ApiEditOrder:             banexg.HasOk,
					banexg.ApiCancelOrder:           banexg.HasOk,
				},
			},
			CredKeys: map[string]bool{"ApiKey": true, "Secret": true},
		},
	}
	
	exg.Sign = makeSign(exg)
	exg.FetchMarkets = makeFetchMarkets(exg)
	err := exg.Init()
	return exg, err
}

func NewExchange(Options map[string]interface{}) (banexg.BanExchange, *errs.Error) {
	return New(Options)
}