package longportapp

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/banbox/banexg"
	"github.com/banbox/banexg/errs"
	"github.com/banbox/banexg/utils"
)

// toString 将任意类型转换为字符串
func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

// makeSign 创建签名函数
func makeSign(e *LongPortApp) banexg.FuncSign {
	return func(api *banexg.Entry, params map[string]interface{}) *banexg.HttpReq {
		var args = utils.SafeParams(params)
		accID := e.PopAccName(args)
		
		// 获取账户凭证
		accID, creds, err := e.GetAccountCreds(accID)
		if err != nil {
			return &banexg.HttpReq{Error: err, Private: true}
		}
		
		// 构建请求URL
		baseURL := e.Hosts.GetHost(api.Host)
		fullURL := baseURL + "/" + api.Path
		
		// 添加查询参数
		if len(args) > 0 && api.Method == "GET" {
			values := url.Values{}
			for k, v := range args {
				values.Add(k, toString(v))
			}
			if len(values) > 0 {
				fullURL += "?" + values.Encode()
			}
		}
		
		// 构建请求头
		headers := http.Header{}
		headers.Set("Authorization", "Bearer "+creds.Password) // AccessToken
		headers.Set("Content-Type", "application/json")
		headers.Set("User-Agent", e.UserAgent)
		
		// 构建请求体
		body := ""
		if api.Method == "POST" || api.Method == "PUT" {
			if len(args) > 0 {
				bodyBytes, err := utils.MarshalString(args)
				if err != nil {
					return &banexg.HttpReq{Error: errs.NewMsg(errs.CodeParamInvalid, "marshal body failed: %v", err)}
				}
				body = bodyBytes
			}
		}
		
		return &banexg.HttpReq{
			AccName: accID,
			Url:     fullURL,
			Method:  api.Method,
			Headers: headers,
			Body:    body,
			Private: true,
		}
	}
}

// makeFetchMarkets 创建获取市场信息的函数
func makeFetchMarkets(e *LongPortApp) banexg.FuncFetchMarkets {
	return func(marketTypes []string, params map[string]interface{}) (banexg.MarketMap, *errs.Error) {
		markets := make(banexg.MarketMap)
		
		// LongPort主要支持港股、美股、A股，这里创建一些示例市场
		// 实际使用时应该通过API获取可交易的股票列表
		
		// 港股示例
		markets["700.HK"] = &banexg.Market{
			ID:       "700.HK",
			Symbol:   "700.HK",
			Base:     "700",
			Quote:    "HKD",
			Type:     banexg.MarketSpot,
			Spot:     true,
			Active:   true,
			Taker:    0.003,
			Maker:    0.003,
			Precision: &banexg.Precision{
				Amount: 0.01,
				Price:  0.01,
			},
			Limits: &banexg.MarketLimits{
				Amount: &banexg.LimitRange{Min: 1, Max: 1000000},
				Price:  &banexg.LimitRange{Min: 0.01, Max: 10000},
			},
		}
		
		// 美股示例
		markets["AAPL.US"] = &banexg.Market{
			ID:       "AAPL.US",
			Symbol:   "AAPL.US",
			Base:     "AAPL",
			Quote:    "USD",
			Type:     banexg.MarketSpot,
			Spot:     true,
			Active:   true,
			Taker:    0.003,
			Maker:    0.003,
			Precision: &banexg.Precision{
				Amount: 0.01,
				Price:  0.01,
			},
			Limits: &banexg.MarketLimits{
				Amount: &banexg.LimitRange{Min: 1, Max: 1000000},
				Price:  &banexg.LimitRange{Min: 0.01, Max: 10000},
			},
		}
		
		// A股示例
		markets["000001.SZ"] = &banexg.Market{
			ID:       "000001.SZ",
			Symbol:   "000001.SZ",
			Base:     "000001",
			Quote:    "CNY",
			Type:     banexg.MarketSpot,
			Spot:     true,
			Active:   true,
			Taker:    0.003,
			Maker:    0.003,
			Precision: &banexg.Precision{
				Amount: 0.01,
				Price:  0.01,
			},
			Limits: &banexg.MarketLimits{
				Amount: &banexg.LimitRange{Min: 100, Max: 1000000},
				Price:  &banexg.LimitRange{Min: 0.01, Max: 10000},
			},
		}
		
		return markets, nil
	}
}