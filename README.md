# Exchange Failover

交易所 API 備援機制 Go 模組，支援自動切換與恢復。

## 功能特色

- 🔄 **自動備援**：主交易所異常時自動切換到備援交易所
- ⏰ **自動恢復**：主交易所恢復後自動切回
- 🔒 **鎖定機制**：防止頻繁切換
- 📊 **錯誤計數**：可配置錯誤閾值與計數有效期
- 🔔 **告警通知**：切換與恢復時發送通知

## 安裝

```bash
go get github.com/yourorg/exchange-failover
```

## 快速開始

```go
package main

import (
    "github.com/yourorg/exchange-failover"
    "github.com/redis/go-redis/v9"
)

func main() {
    // 1. 初始化 Redis
    rdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

    // 2. 建立 Connector (你需要實作 ExchangeConnector 介面)
    binance := NewBinanceConnector()
    okx := NewOKXConnector()

    // 3. 建立 Proxy
    proxy := failover.NewProxy(
        failover.WithPrimaryConnector(binance),
        failover.WithStandbyConnector(okx),
        failover.WithCache(rdb),
        failover.WithAlertService(myAlertService),
    )

    // 4. 建立 Adapter
    api := failover.NewAdapter(proxy)

    // 5. 使用 API
    price, err := api.NewestQuoteTicker("BTCUSDT")
}
```

## 架構

```
┌─────────────┐
│  業務代碼    │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ ExchangeApi │  (高層 API)
└──────┬──────┘
       │
       ▼
┌─────────────────────┐
│ ExchangeApiProxy    │  (備援邏輯)
└──────┬──────────────┘
       │
       ▼
┌──────────────┐   ┌─────────────┐
│  Binance    │   │    OKX      │
│ (Primary)   │   │ (Standby)   │
└──────────────┘   └─────────────┘
```

## 介面說明

### ExchangeConnector

交易所連接器介面，定義交易所 API 操作。

```go
type ExchangeConnector interface {
    IsSystemAbnormal(FailureCode string) bool
    NewestQuoteTicker(symbol string) (res ExchangeApiResponse, err error)
    SpotTrade(symbol, side, quantity, price string) (res ExchangeApiResponse, err error)
    // ... 其他方法
}
```

### ExchangeApiProxy

Proxy 介面，負責備援機制。

```go
type ExchangeApiProxy interface {
    Invoke(
        fn func(ct ExchangeConnectorType, connector ExchangeConnector) (ExchangeApiResponse, error),
        con *ExchangeConnectorType,
        needStandbyConnector bool,
    ) (ExchangeApiResponse, error)
    NowConnect() string
}
```

### ExchangeApi

高層 API 介面，回傳解析後的業務數據。

```go
type ExchangeApi interface {
    NewestQuoteTicker(symbol string) (price decimal.Decimal, err error)
    SpotTrade(symbol, side, quantity, price string) (output map[string]interface{}, err error)
    // ... 其他方法
}
```

## 設定

```go
config := failover.Config{
    ErrThreshold:      5,           // 錯誤次數閾值
    ErrTTL:            30 * time.Second,  // 錯誤計數有效期
    LockTimeTTL:       30 * time.Minute, // LockTime 有效期
}

proxy := failover.NewProxy(
    failover.WithConfig(config),
    // ... 其他選項
)
```

## 備援機制

### 觸發條件
- 30 秒內發生 5 次系統異常 → 切換到 OKX

### 恢復條件
- LockTime (30 分鐘) 過期後
- 嘗試切回主交易所
- API 調用成功 → 切回主交易所

##  Licence

MIT
