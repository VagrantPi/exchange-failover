package failover

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
)

const ExchangeConnectorKey = "exchange:connector"
const ExchangeConnectorLockTime = "exchange:lockTime"
const ExchangeConnectorLockTimeTTL = time.Duration(30) * time.Minute
const ExchangeConnectorErrTimeAt = "exchange:errTime"
const ExchangeConnectorErrThreshold = 5
const ExchangeConnectorErrTTL = time.Duration(30) * time.Second

type ExchangeApiProxyImpl struct {
	BinanceImpl  ExchangeConnector
	OKXImpl      ExchangeConnector
	Cache        redis.UniversalClient
	AlertService IAlertService
}

func (proxy ExchangeApiProxyImpl) getConnector(con *ExchangeConnectorType, needStandbyConnector bool) (ct ExchangeConnectorType, connector ExchangeConnector, err error) {
	ctx := context.Background()

	if con != nil {
		switch *con {
		case ExchangeConnectorTypeBinance:
			return ExchangeConnectorTypeBinance, proxy.BinanceImpl, nil
		case ExchangeConnectorTypeOKX:
			return ExchangeConnectorTypeOKX, proxy.OKXImpl, nil
		}
	}

	nowConnector, err := proxy.Cache.Get(ctx, ExchangeConnectorKey).Result()
	if err != nil && err != redis.Nil {
		return "", nil, err
	}
	if nowConnector == ExchangeConnectorTypeBinance.String() || nowConnector == "" {
		return ExchangeConnectorTypeBinance, proxy.BinanceImpl, nil
	}

	exist, err := proxy.Cache.Exists(ctx, ExchangeConnectorLockTime).Result()
	if err != nil {
		return "", nil, err
	}
	if exist == 0 && needStandbyConnector {
		return ExchangeConnectorTypeBinance, proxy.BinanceImpl, nil
	}

	return ExchangeConnectorTypeOKX, proxy.OKXImpl, nil
}

func (proxy ExchangeApiProxyImpl) addFailureCount(ct ExchangeConnectorType) error {
	ctx := context.Background()

	nowConnector, err := proxy.Cache.Get(ctx, ExchangeConnectorKey).Result()
	if err != nil && err != redis.Nil {
		return err
	}
	log.Infof("addFailureCount nowConnector: %v", nowConnector)
	if nowConnector == ExchangeConnectorTypeOKX.String() {
		err = proxy.Cache.Set(ctx, ExchangeConnectorLockTime, time.Now(), ExchangeConnectorLockTimeTTL).Err()
		if err != nil {
			return err
		}
	}

	key := fmt.Sprintf("%v:%v:%v", ExchangeConnectorErrTimeAt, ct, time.Now().UnixMilli())
	if err := proxy.Cache.Set(ctx, key, true, ExchangeConnectorErrTTL).Err(); err != nil {
		return err
	}

	key = fmt.Sprintf("%v:%v:*", ExchangeConnectorErrTimeAt, ct)
	errTimestamps, err := proxy.Cache.Keys(ctx, key).Result()
	if err != nil {
		return err
	}

	if len(errTimestamps) >= ExchangeConnectorErrThreshold {
		err = proxy.Cache.Set(ctx, ExchangeConnectorKey, ExchangeConnectorTypeOKX.String(), -1).Err()
		if err != nil {
			return err
		}

		err = proxy.Cache.Set(ctx, ExchangeConnectorLockTime, time.Now(), ExchangeConnectorLockTimeTTL).Err()
		if err != nil {
			return err
		}

		msg := "因 幣安 發生異常無法使用，先採用 OKX 進行避險、報價的執行。請通知第三方廠商做緊急處理。"
		if innerErr := proxy.AlertService.SendErrorAlert("Binance", msg); innerErr != nil {
			log.Infof("HandleRequestAlert error: %v", innerErr)
		}
	}

	return nil
}

func (proxy ExchangeApiProxyImpl) resetFailureCount(ct ExchangeConnectorType) error {
	ctx := context.Background()

	existConnectorLockTime, err := proxy.Cache.Exists(ctx, ExchangeConnectorLockTime).Result()
	if err != nil {
		return nil
	}
	isLockOKX := existConnectorLockTime != 0

	nowConnector, err := proxy.Cache.Get(ctx, ExchangeConnectorKey).Result()
	if err != nil && err != redis.Nil {
		return err
	}

	if nowConnector == ExchangeConnectorTypeBinance.String() || isLockOKX {
		return nil
	}

	key := fmt.Sprintf("%v:%v:*", ExchangeConnectorErrTimeAt, ct)
	errTimestamps, err := proxy.Cache.Keys(ctx, key).Result()
	if err != nil {
		return err
	}
	for _, k := range errTimestamps {
		err = proxy.Cache.Del(ctx, k).Err()
		if err != nil {
			return err
		}
	}

	err = proxy.Cache.Set(context.Background(), ExchangeConnectorKey, ExchangeConnectorTypeBinance.String(), -1).Err()
	if err != nil {
		return err
	}

	if innerErr := proxy.AlertService.SendRecoveryAlert("Binance"); innerErr != nil {
		log.Infof("HandleRequestAlert error: %v", innerErr)
	}
	return nil
}

func (proxy ExchangeApiProxyImpl) NowConnect() string {
	nowConnector, err := proxy.Cache.Get(context.Background(), ExchangeConnectorKey).Result()
	if err != nil && err != redis.Nil {
		return ExchangeConnectorTypeBinance.String()
	}

	return nowConnector
}

func (proxy ExchangeApiProxyImpl) Invoke(fn func(ct ExchangeConnectorType, connector ExchangeConnector) (ExchangeApiResponse, error), con *ExchangeConnectorType, needStandbyConnector bool) (ExchangeApiResponse, error) {
	cType, connector, err := proxy.getConnector(con, needStandbyConnector)
	if err != nil {
		return ExchangeApiResponse{}, fmt.Errorf("getConnector error: %w", err)
	}

	apiResponse, err := fn(cType, connector)
	if apiResponse.IsSuccess {
		err = proxy.resetFailureCount(cType)
		if err != nil {
			return ExchangeApiResponse{}, fmt.Errorf("reset failure count err: %w", err)
		}
	}
	if !apiResponse.IsSuccess {
		if connector.IsSystemAbnormal(apiResponse.FailureCode) {
			err = proxy.addFailureCount(cType)
			if err != nil {
				return ExchangeApiResponse{}, fmt.Errorf("add failure count err: %w", err)
			}
		}

		return ExchangeApiResponse{},
			fmt.Errorf("call api error: IsSuccess=%v, Body=%v, FailureCode=%v, ConnectorType=%v",
				apiResponse.IsSuccess, string(apiResponse.Body), apiResponse.FailureCode, apiResponse.ConnectorType)

	}

	return apiResponse, nil
}
