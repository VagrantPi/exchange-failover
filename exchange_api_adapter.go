package failover

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

type ExchangeApiAdapter struct {
	ApiProxy ExchangeApiProxy
}

func (e ExchangeApiAdapter) NowConnect() string {
	return e.ApiProxy.NowConnect()
}

func (e ExchangeApiAdapter) Klines(symbol string, interval string, limit uint64) (klines []map[string]interface{}, err error) {
	apiResponse, err := e.ApiProxy.Invoke(func(cType ExchangeConnectorType, connector ExchangeConnector) (ExchangeApiResponse, error) {
		return connector.Klines(symbol, interval, limit)
	}, nil, false)
	if err != nil {
		return nil, err
	}

	klines = []map[string]interface{}{}
	err = json.Unmarshal(apiResponse.Body, &klines)
	if err != nil {
		return nil, err
	}

	return klines, nil
}

func (e ExchangeApiAdapter) ClosingTimeRemaining(interval string) time.Duration {
	apiResponse, err := e.ApiProxy.Invoke(func(cType ExchangeConnectorType, connector ExchangeConnector) (ExchangeApiResponse, error) {
		return connector.ClosingTimeRemaining(interval)
	}, nil, false)

	remainingTime := time.Duration(0)
	err = json.Unmarshal(apiResponse.Body, &remainingTime)
	if err != nil {
		return remainingTime
	}

	return remainingTime
}

func (e ExchangeApiAdapter) GetPriceHistoryIntervalLimit(intervalLetter string) (interval string, limit uint64) {
	apiResponse, err := e.ApiProxy.Invoke(func(cType ExchangeConnectorType, connector ExchangeConnector) (ExchangeApiResponse, error) {
		return connector.GetPriceHistoryIntervalLimit(intervalLetter)
	}, nil, false)
	if err != nil {
		return "", 0
	}

	resp := struct {
		Interval string `json:"interval"`
		Limit    uint64 `json:"limit"`
	}{}
	err = json.Unmarshal(apiResponse.Body, &resp)
	if err != nil {
		return "", 0
	}

	return resp.Interval, resp.Limit
}

func (e ExchangeApiAdapter) FutureTrade(symbol, side, quantity, price string) (output map[string]interface{}, err error) {
	apiResponse, err := e.ApiProxy.Invoke(func(cType ExchangeConnectorType, connector ExchangeConnector) (ExchangeApiResponse, error) {
		return connector.FutureTrade(symbol, side, quantity, price)
	}, nil, true)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{}
	err = json.Unmarshal(apiResponse.Body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (e ExchangeApiAdapter) GetUSDTMFuturesPrecision(base string) (pricePrecision, quantityPrecision int32, err error) {
	apiResponse, err := e.ApiProxy.Invoke(func(cType ExchangeConnectorType, connector ExchangeConnector) (ExchangeApiResponse, error) {
		return connector.GetUSDTMFuturesPrecision(base)
	}, nil, false)
	if err != nil {
		return 0, 0, err
	}

	result := struct {
		PricePrecision    int32 `json:"pricePrecision"`
		QuantityPrecision int32 `json:"quantityPrecision"`
	}{}
	err = json.Unmarshal(apiResponse.Body, &result)
	if err != nil {
		return 0, 0, err
	}

	return result.PricePrecision, result.QuantityPrecision, nil
}

func (e ExchangeApiAdapter) SpotTrade(symbol, side, quantity, price string) (output map[string]interface{}, err error) {
	apiResponse, err := e.ApiProxy.Invoke(func(cType ExchangeConnectorType, connector ExchangeConnector) (ExchangeApiResponse, error) {
		return connector.SpotTrade(symbol, side, quantity, price)
	}, nil, false)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{}
	err = json.Unmarshal(apiResponse.Body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (e ExchangeApiAdapter) FuturesExchangeInfo(symbol string) (resp map[string]interface{}, err error) {
	apiResponse, err := e.ApiProxy.Invoke(func(cType ExchangeConnectorType, connector ExchangeConnector) (ExchangeApiResponse, error) {
		return connector.FuturesExchangeInfo(symbol)
	}, nil, false)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{}
	err = json.Unmarshal(apiResponse.Body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (e ExchangeApiAdapter) GetFuturesBills(startTime int64) (resp []map[string]interface{}, err error) {
	apiResponse, err := e.ApiProxy.Invoke(func(cType ExchangeConnectorType, connector ExchangeConnector) (ExchangeApiResponse, error) {
		return connector.GetFuturesBills(startTime)
	}, nil, false)
	if err != nil {
		return []map[string]interface{}{}, err
	}

	result := []map[string]interface{}{}
	err = json.Unmarshal(apiResponse.Body, &result)
	if err != nil {
		return []map[string]interface{}{}, err
	}

	return result, nil
}

func (e ExchangeApiAdapter) FuturesTransfer(symbol, amount, transferType string, connector ExchangeConnectorType) (err error) {
	apiResponse, err := e.ApiProxy.Invoke(func(cType ExchangeConnectorType, connector ExchangeConnector) (ExchangeApiResponse, error) {
		return connector.FuturesTransfer(symbol, amount, transferType)
	}, &connector, false)
	if err != nil {
		return err
	}

	parseErr := json.Unmarshal(apiResponse.Body, &err)
	if parseErr != nil {
		return parseErr
	}

	return parseErr
}

func (e ExchangeApiAdapter) FuturesAccount() (account map[string]interface{}, err error) {
	apiResponse, err := e.ApiProxy.Invoke(func(cType ExchangeConnectorType, connector ExchangeConnector) (ExchangeApiResponse, error) {
		return connector.FuturesAccount()
	}, nil, false)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{}
	err = json.Unmarshal(apiResponse.Body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (e ExchangeApiAdapter) FuturesAccountPositionRisk(symbol string) (risk []map[string]interface{}, err error) {
	apiResponse, err := e.ApiProxy.Invoke(func(cType ExchangeConnectorType, connector ExchangeConnector) (ExchangeApiResponse, error) {
		return connector.FuturesAccountPositionRisk(symbol)
	}, nil, false)
	if err != nil {
		return []map[string]interface{}{}, err
	}

	result := []map[string]interface{}{}
	err = json.Unmarshal(apiResponse.Body, &result)
	if err != nil {
		return []map[string]interface{}{}, err
	}

	return result, nil
}

func (e ExchangeApiAdapter) SpotAllOrders(symbol string, limit int64) (output []map[string]interface{}, err error) {
	binanceCon := ExchangeConnectorTypeBinance

	apiResponse, err := e.ApiProxy.Invoke(func(cType ExchangeConnectorType, connector ExchangeConnector) (ExchangeApiResponse, error) {
		return connector.SpotAllOrders(symbol, limit)
	}, &binanceCon, false)
	if err != nil {
		return []map[string]interface{}{}, err
	}

	result := []map[string]interface{}{}
	err = json.Unmarshal(apiResponse.Body, &result)
	if err != nil {
		return []map[string]interface{}{}, err
	}

	return result, nil
}

func (e ExchangeApiAdapter) SpotAccountTradeList(symbol string, limit int64) (output []map[string]interface{}, err error) {
	binanceCon := ExchangeConnectorTypeBinance

	apiResponse, err := e.ApiProxy.Invoke(func(cType ExchangeConnectorType, connector ExchangeConnector) (ExchangeApiResponse, error) {
		return connector.SpotAccountTradeList(symbol, limit)
	}, &binanceCon, false)
	if err != nil {
		return []map[string]interface{}{}, err
	}

	result := []map[string]interface{}{}
	err = json.Unmarshal(apiResponse.Body, &result)
	if err != nil {
		return []map[string]interface{}{}, err
	}

	return result, nil
}

func (e ExchangeApiAdapter) GetCommission(symbols string) (output []map[string]interface{}, err error) {
	apiResponse, err := e.ApiProxy.Invoke(func(cType ExchangeConnectorType, connector ExchangeConnector) (ExchangeApiResponse, error) {
		return connector.GetCommission(symbols)
	}, nil, false)
	if err != nil {
		return []map[string]interface{}{}, err
	}

	result := []map[string]interface{}{}
	err = json.Unmarshal(apiResponse.Body, &result)
	if err != nil {
		return []map[string]interface{}{}, err
	}

	return result, nil
}

func (e ExchangeApiAdapter) PerpAccountTradeList(symbol string, limit int64) (output []map[string]interface{}, err error) {
	binanceCon := ExchangeConnectorTypeBinance

	apiResponse, err := e.ApiProxy.Invoke(func(cType ExchangeConnectorType, connector ExchangeConnector) (ExchangeApiResponse, error) {
		return connector.PerpAccountTradeList(symbol, limit)
	}, &binanceCon, false)
	if err != nil {
		return []map[string]interface{}{}, err
	}

	result := []map[string]interface{}{}
	err = json.Unmarshal(apiResponse.Body, &result)
	if err != nil {
		return []map[string]interface{}{}, err
	}

	return result, nil
}

func (e ExchangeApiAdapter) SpotAccountInternalTransferRecord(startTime, endTime int64) (output []map[string]interface{}, err error) {
	binanceCon := ExchangeConnectorTypeBinance

	apiResponse, err := e.ApiProxy.Invoke(func(cType ExchangeConnectorType, connector ExchangeConnector) (ExchangeApiResponse, error) {
		return connector.SpotAccountInternalTransferRecord(startTime, endTime)
	}, &binanceCon, false)
	if err != nil {
		return []map[string]interface{}{}, err
	}

	result := []map[string]interface{}{}
	err = json.Unmarshal(apiResponse.Body, &result)
	if err != nil {
		return []map[string]interface{}{}, err
	}

	return result, nil
}

func (e ExchangeApiAdapter) SpotWithdraw(symbol, amount, to, network string) (id string, err error) {
	binanceCon := ExchangeConnectorTypeBinance

	apiResponse, err := e.ApiProxy.Invoke(func(cType ExchangeConnectorType, connector ExchangeConnector) (ExchangeApiResponse, error) {
		return connector.SpotWithdraw(symbol, amount, to, network)
	}, &binanceCon, false)
	if err != nil {
		return "", err
	}

	result := struct {
		Id string `json:"id"`
	}{}
	err = json.Unmarshal(apiResponse.Body, &result)
	if err != nil {
		return "", err
	}

	return result.Id, nil
}

func (e ExchangeApiAdapter) SpotWithdrawRecord(startTime, endTime int64) (output []map[string]interface{}, err error) {
	binanceCon := ExchangeConnectorTypeBinance

	apiResponse, err := e.ApiProxy.Invoke(func(cType ExchangeConnectorType, connector ExchangeConnector) (ExchangeApiResponse, error) {
		return connector.SpotWithdrawRecord(startTime, endTime)
	}, &binanceCon, false)
	if err != nil {
		return []map[string]interface{}{}, err
	}

	result := []map[string]interface{}{}
	err = json.Unmarshal(apiResponse.Body, &result)
	if err != nil {
		return []map[string]interface{}{}, err
	}

	return result, nil
}

func (e ExchangeApiAdapter) CapitalCoinGetAll() (coinConfigs []map[string]interface{}, err error) {
	binanceCon := ExchangeConnectorTypeBinance

	apiResponse, err := e.ApiProxy.Invoke(func(cType ExchangeConnectorType, connector ExchangeConnector) (ExchangeApiResponse, error) {
		return connector.CapitalCoinGetAll()
	}, &binanceCon, false)
	if err != nil {
		return []map[string]interface{}{}, err
	}

	result := []map[string]interface{}{}
	err = json.Unmarshal(apiResponse.Body, &result)
	if err != nil {
		return []map[string]interface{}{}, err
	}

	return result, nil
}

func (e ExchangeApiAdapter) SpotAssets(symbol string) (spotAssets []map[string]interface{}, err error) {
	apiResponse, err := e.ApiProxy.Invoke(func(cType ExchangeConnectorType, connector ExchangeConnector) (ExchangeApiResponse, error) {
		return connector.SpotAssets(symbol)
	}, nil, false)
	if err != nil {
		return []map[string]interface{}{}, err
	}

	result := []map[string]interface{}{}
	err = json.Unmarshal(apiResponse.Body, &result)
	if err != nil {
		return []map[string]interface{}{}, err
	}

	return result, nil
}

func (e ExchangeApiAdapter) NewestQuoteTicker(symbol string) (price decimal.Decimal, err error) {
	apiResponse, err := e.ApiProxy.Invoke(func(cType ExchangeConnectorType, connector ExchangeConnector) (ExchangeApiResponse, error) {
		return connector.NewestQuoteTicker(symbol)
	}, nil, false)
	if err != nil {
		return decimal.Zero, err
	}

	result := struct {
		Price string `json:"price"`
	}{}
	err = json.Unmarshal(apiResponse.Body, &result)
	if err != nil {
		return decimal.Zero, err
	}

	d, err := decimal.NewFromString(result.Price)
	if err != nil {
		return decimal.Zero, err
	}

	return d, nil
}

func (e ExchangeApiAdapter) GetSpotPrecision(base string) (pricePrecision int32, quantityPrecision int32, quoteQuantityPrecision int32, err error) {
	apiResponse, err := e.ApiProxy.Invoke(func(cType ExchangeConnectorType, connector ExchangeConnector) (ExchangeApiResponse, error) {
		return connector.GetSpotPrecision(base)
	}, nil, false)
	if err != nil {
		return 0, 0, 0, err
	}

	result := struct {
		PricePrecision          int32 `json:"pricePrecision"`
		QuantityPrecision       int32 `json:"quantityPrecision"`
		QuoteQuantityPrecision int32 `json:"quoteQuantityPrecision"`
	}{}
	err = json.Unmarshal(apiResponse.Body, &result)
	if err != nil {
		return 0, 0, 0, err
	}

	return result.PricePrecision, result.QuantityPrecision, result.QuoteQuantityPrecision, nil
}

func (e ExchangeApiAdapter) SymbolPriceTicker() (price []map[string]interface{}, err error) {
	apiResponse, err := e.ApiProxy.Invoke(func(cType ExchangeConnectorType, connector ExchangeConnector) (ExchangeApiResponse, error) {
		return connector.SymbolPriceTicker()
	}, nil, false)
	if err != nil {
		return []map[string]interface{}{}, err
	}

	result := []map[string]interface{}{}
	err = json.Unmarshal(apiResponse.Body, &result)
	if err != nil {
		return []map[string]interface{}{}, err
	}

	return result, nil
}
