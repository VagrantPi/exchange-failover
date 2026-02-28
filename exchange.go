package failover

import (
	"time"

	"github.com/shopspring/decimal"
)

type ExchangeConnectorType string

const (
	ExchangeConnectorTypeBinance ExchangeConnectorType = "Binance"
	ExchangeConnectorTypeOKX     ExchangeConnectorType = "OKX"
)

func (ect ExchangeConnectorType) String() string {
	return string(ect)
}

type ExchangeApiResponse struct {
	IsSuccess     bool
	Body          []byte
	FailureCode   string
	ConnectorType ExchangeConnectorType
}

type ExchangeConnector interface {
	IsSystemAbnormal(FailureCode string) bool
	Klines(symbol string, interval string, limit uint64) (res ExchangeApiResponse, err error)
	ClosingTimeRemaining(interval string) (res ExchangeApiResponse, err error)
	GetPriceHistoryIntervalLimit(intervalLetter string) (res ExchangeApiResponse, err error)
	FutureTrade(symbol, side, quantity, price string) (res ExchangeApiResponse, err error)
	GetUSDTMFuturesPrecision(base string) (res ExchangeApiResponse, err error)
	SpotTrade(symbol, side, quantity, price string) (res ExchangeApiResponse, err error)
	FuturesExchangeInfo(symbol string) (res ExchangeApiResponse, err error)
	GetFuturesBills(startTime int64) (res ExchangeApiResponse, err error)
	FuturesTransfer(symbol, amount, transferType string) (res ExchangeApiResponse, err error)
	FuturesAccount() (res ExchangeApiResponse, err error)
	FuturesAccountPositionRisk(symbol string) (res ExchangeApiResponse, err error)
	SpotAllOrders(symbol string, limit int64) (res ExchangeApiResponse, err error)
	SpotAccountTradeList(symbol string, limit int64) (res ExchangeApiResponse, err error)
	PerpAccountTradeList(symbol string, limit int64) (res ExchangeApiResponse, err error)
	GetCommission(symbols string) (res ExchangeApiResponse, err error)
	SpotAccountInternalTransferRecord(startTime, endTime int64) (res ExchangeApiResponse, err error)
	SpotWithdraw(symbol, amount, to, network string) (res ExchangeApiResponse, err error)
	SpotWithdrawRecord(startTime, endTime int64) (res ExchangeApiResponse, err error)
	CapitalCoinGetAll() (res ExchangeApiResponse, err error)
	SpotAssets(symbol string) (res ExchangeApiResponse, err error)
	NewestQuoteTicker(symbol string) (res ExchangeApiResponse, err error)
	GetSpotPrecision(base string) (res ExchangeApiResponse, err error)
	SymbolPriceTicker() (res ExchangeApiResponse, err error)
}

type ExchangeApiProxy interface {
	Invoke(fn func(ct ExchangeConnectorType, connector ExchangeConnector) (ExchangeApiResponse, error), con *ExchangeConnectorType, needStandbyConnector bool) (ExchangeApiResponse, error)
	NowConnect() string
}

type ExchangeApi interface {
	NowConnect() string
	Klines(symbol string, interval string, limit uint64) (klines []map[string]interface{}, err error)
	ClosingTimeRemaining(interval string) time.Duration
	GetPriceHistoryIntervalLimit(intervalLetter string) (interval string, limit uint64)
	FutureTrade(symbol, side, quantity, price string) (output map[string]interface{}, err error)
	GetUSDTMFuturesPrecision(base string) (pricePrecision, quantityPrecision int32, err error)
	SpotTrade(symbol, side, quantity, price string) (output map[string]interface{}, err error)
	FuturesExchangeInfo(symbol string) (resp map[string]interface{}, err error)
	GetFuturesBills(startTime int64) (resp []map[string]interface{}, err error)
	FuturesTransfer(symbol, amount, transferType string, connector ExchangeConnectorType) (err error)
	FuturesAccount() (account map[string]interface{}, err error)
	FuturesAccountPositionRisk(symbol string) (risk []map[string]interface{}, err error)
	SpotAllOrders(symbol string, limit int64) (output []map[string]interface{}, err error)
	SpotAccountTradeList(symbol string, limit int64) (output []map[string]interface{}, err error)
	PerpAccountTradeList(symbol string, limit int64) (output []map[string]interface{}, err error)
	GetCommission(symbols string) (output []map[string]interface{}, err error)
	SpotAccountInternalTransferRecord(startTime, endTime int64) (output []map[string]interface{}, err error)
	SpotWithdraw(symbol, amount, to, network string) (id string, err error)
	SpotWithdrawRecord(startTime, endTime int64) (output []map[string]interface{}, err error)
	CapitalCoinGetAll() (coinConfigs []map[string]interface{}, err error)
	SpotAssets(symbol string) (spotAssets []map[string]interface{}, err error)
	NewestQuoteTicker(symbol string) (price decimal.Decimal, err error)
	GetSpotPrecision(base string) (pricePrecision int32, quantityPrecision int32, quoteQuantityPrecision int32, err error)
	SymbolPriceTicker() (price []map[string]interface{}, err error)
}

type IAlertService interface {
	SendErrorAlert(source, msg string) error
	SendRecoveryAlert(source string) error
}
