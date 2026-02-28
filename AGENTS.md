# AGENTS.md - Exchange Failover 專案開發指南

本檔案提供給 Agentic 編碼代理人在此專案中開發時的參考指引。

## 專案概述

- **語言**: Go 1.20+
- **框架**: Kratos
- **主要依賴**: 
  - `github.com/go-kratos/kratos/v2` - 日誌
  - `github.com/redis/go-redis/v9` - 快取
  - `github.com/shopspring/decimal` - 精確小數運算
- **用途**: 交易所 API 備援機制，支援 Binance/OKX 自動切換與恢復

---

## 建置與測試指令

### 建置
```bash
# 編譯專案
go build ./...

# 下載依賴
go mod download
go mod tidy
```

### 測試
```bash
# 執行所有測試
go test ./...

# 執行單一測試
go test -v -run TestFunctionName ./...

# 顯示測試覆蓋率
go test -cover ./...
```

### 程式碼品質
```bash
# Vet 靜態分析
go vet ./...

# Fmt 格式檢查
go fmt ./...

# Lint (需安裝 golint)
golangci-lint run ./...
```

### 整合開發
```bash
# 產生 mock code (若使用 go generate)
go generate ./...
```

---

## 程式碼風格指南

### 命名慣例

- **套件名稱**: 使用簡短小寫，如 `failover`
- **匯出函式/類型**: PascalCase，如 `NewProxy`, `ExchangeApiProxyImpl`
- **未匯出成員**: camelCase，如 `binanceImpl`, `apiProxy`
- **常數**: 全大寫 + 底線，如 `ExchangeConnectorKey`
- **介面命名**: 以 `I` 開頭或遵循 Go 慣例（如 `ExchangeConnector`）

### 檔案組織

- 一個套件一個目錄（除非測試檔案）
- 檔案名稱使用小寫底線分隔：`exchange_proxy.go`, `exchange_api_adapter.go`
- 測試檔案命名：`xxx_test.go`

### 匯入排序

標準庫優先，三方庫次之，空白行分隔：

```go
import (
    "context"
    "encoding/json"
    "time"

    "github.com/go-kratos/kratos/v2/log"
    "github.com/redis/go-redis/v9"
    "github.com/shopspring/decimal"
)
```

### 錯誤處理

- 錯誤作為最後回傳值：`func foo() (result Type, err error)`
- 使用 `fmt.Errorf` 搭配 `%w` 包裝錯誤：`fmt.Errorf("getConnector error: %w", err)`
- 避免 naked return，明確回傳錯誤
- Redis 操作需檢查 `redis.Nil` 情況

### 類型與介面

- 優先使用介面定義行為（如 `ExchangeConnector`, `ExchangeApiProxy`）
- 使用 Options Pattern 進行可選參數配置
- 結構體標籤用於 JSON 序列化：`json:"pricePrecision"`

### 日誌與監控

- 使用 Kratos log 套件：`log.Infof`, `log.Errorf`
- 避免在生產環境使用 `fmt.Println`

### 註冊與依賴注入

- 透過建構函式或 Options 注入依賴（如 `redis.UniversalClient`, `IAlertService`）
- 避免全域變數狀態

---

## 架構說明

### 核心介面

```
ExchangeConnector     → 交易所連接器 (Binance/OKX)
    │
    ▼
ExchangeApiProxy      → 備援邏輯 (自動切換/恢復)
    │
    ▼
ExchangeApi           → 高層 API (業務使用)
```

### 新增功能流程

1. 在 `ExchangeConnector` 介面新增方法
2. 在 `ExchangeApiProxyImpl` 實作 proxy 邏輯
3. 在 `ExchangeApiAdapter` 封裝並解析 response

---

## 常見任務參考

### 新增交易所支援

1. 在 `exchange.go` 新增 `ExchangeConnectorType` 常數
2. 實作 `ExchangeConnector` 介面
3. 在 `NewProxy` 註冊新 connector

### 修改備援閾值

```go
proxy := failover.NewProxy(
    failover.WithConfig(failover.Config{
        ErrThreshold: 10,      // 錯誤次數閾值
        ErrTTL: 60 * time.Second,    // 錯誤計數有效期
        LockTimeTTL: 60 * time.Minute, // LockTime 有效期
    }),
    // ...
)
```

---

## 外部資源

- [Go 官方文檔](https://go.dev/doc/)
- [Kratos Framework](https://go-kratos.dev/)
- [Redis Go Client](https://redis.uptrace.dev/)
- [Decimal 套件](https://github.com/shopspring/decimal)
