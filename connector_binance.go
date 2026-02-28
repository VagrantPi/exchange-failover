package failover

import "fmt"

type BinanceConnector struct {
	apiKey     string
	secretKey  string
	baseURL    string
}

func NewBinanceConnector(apiKey, secretKey, baseURL string) *BinanceConnector {
	return &BinanceConnector{
		apiKey:    apiKey,
		secretKey: secretKey,
		baseURL:   baseURL,
	}
}

func (b *BinanceConnector) IsSystemAbnormal(failureCode string) bool {
	systemAbnormalCodes := []string{
		"-1000", "-1001", "-1002", "-1003", "-1004", "-1005", "-1006", "-1007", "-1008",
		"-1010", "-1011", "-1012", "-1013", "-1014", "-1015", "-1016", "-1020", "-1021", "-1022",
		"-1102", "-1111", "-1121", "-1136",
	}
	for _, code := range systemAbnormalCodes {
		if failureCode == code {
			return true
		}
	}
	return false
}

func (b *BinanceConnector) Klines(symbol, interval string, limit uint64) (ExchangeApiResponse, error) {
	// TODO: 實作 Binance K線 API
	return ExchangeApiResponse{}, fmt.Errorf("not implemented")
}

func (b *BinanceConnector) ClosingTimeRemaining(interval string) (ExchangeApiResponse, error) {
	return ExchangeApiResponse{}, fmt.Errorf("not implemented")
}

func (b *BinanceConnector) GetPriceHistoryIntervalLimit(intervalLetter string) (ExchangeApiResponse, error) {
	return ExchangeApiResponse{}, fmt.Errorf("not implemented")
}

func (b *BinanceConnector) FutureTrade(symbol, side, quantity, price string) (ExchangeApiResponse, error) {
	return ExchangeApiResponse{}, fmt.Errorf("not implemented")
}

func (b *BinanceConnector) GetUSDTMFuturesPrecision(base string) (ExchangeApiResponse, error) {
	return ExchangeApiResponse{}, fmt.Errorf("not implemented")
}

func (b *BinanceConnector) SpotTrade(symbol, side, quantity, price string) (ExchangeApiResponse, error) {
	return ExchangeApiResponse{}, fmt.Errorf("not implemented")
}

func (b *BinanceConnector) FuturesExchangeInfo(symbol string) (ExchangeApiResponse, error) {
	return ExchangeApiResponse{}, fmt.Errorf("not implemented")
}

func (b *BinanceConnector) GetFuturesBills(startTime int64) (ExchangeApiResponse, error) {
	return ExchangeApiResponse{}, fmt.Errorf("not implemented")
}

func (b *BinanceConnector) FuturesTransfer(symbol, amount, transferType string) (ExchangeApiResponse, error) {
	return ExchangeApiResponse{}, fmt.Errorf("not implemented")
}

func (b *BinanceConnector) FuturesAccount() (ExchangeApiResponse, error) {
	return ExchangeApiResponse{}, fmt.Errorf("not implemented")
}

func (b *BinanceConnector) FuturesAccountPositionRisk(symbol string) (ExchangeApiResponse, error) {
	return ExchangeApiResponse{}, fmt.Errorf("not implemented")
}

func (b *BinanceConnector) SpotAllOrders(symbol string, limit int64) (ExchangeApiResponse, error) {
	return ExchangeApiResponse{}, fmt.Errorf("not implemented")
}

func (b *BinanceConnector) SpotAccountTradeList(symbol string, limit int64) (ExchangeApiResponse, error) {
	return ExchangeApiResponse{}, fmt.Errorf("not implemented")
}

func (b *BinanceConnector) PerpAccountTradeList(symbol string, limit int64) (ExchangeApiResponse, error) {
	return ExchangeApiResponse{}, fmt.Errorf("not implemented")
}

func (b *BinanceConnector) GetCommission(symbols string) (ExchangeApiResponse, error) {
	return ExchangeApiResponse{}, fmt.Errorf("not implemented")
}

func (b *BinanceConnector) SpotAccountInternalTransferRecord(startTime, endTime int64) (ExchangeApiResponse, error) {
	return ExchangeApiResponse{}, fmt.Errorf("not implemented")
}

func (b *BinanceConnector) SpotWithdraw(symbol, amount, to, network string) (ExchangeApiResponse, error) {
	return ExchangeApiResponse{}, fmt.Errorf("not implemented")
}

func (b *BinanceConnector) SpotWithdrawRecord(startTime, endTime int64) (ExchangeApiResponse, error) {
	return ExchangeApiResponse{}, fmt.Errorf("not implemented")
}

func (b *BinanceConnector) CapitalCoinGetAll() (ExchangeApiResponse, error) {
	return ExchangeApiResponse{}, fmt.Errorf("not implemented")
}

func (b *BinanceConnector) SpotAssets(symbol string) (ExchangeApiResponse, error) {
	return ExchangeApiResponse{}, fmt.Errorf("not implemented")
}

func (b *BinanceConnector) NewestQuoteTicker(symbol string) (ExchangeApiResponse, error) {
	// TODO: 實作 Binance 報價 API
	return ExchangeApiResponse{}, fmt.Errorf("not implemented")
}

func (b *BinanceConnector) GetSpotPrecision(base string) (ExchangeApiResponse, error) {
	return ExchangeApiResponse{}, fmt.Errorf("not implemented")
}

func (b *BinanceConnector) SymbolPriceTicker() (ExchangeApiResponse, error) {
	return ExchangeApiResponse{}, fmt.Errorf("not implemented")
}
