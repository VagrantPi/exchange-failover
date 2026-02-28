package failover

import (
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	ErrThreshold       int
	ErrTTL             time.Duration
	LockTimeTTL        time.Duration
	RedisKeyConnector  string
	RedisKeyLockTime   string
	RedisKeyErrTimeAt  string
}

var DefaultConfig = Config{
	ErrThreshold:      5,
	ErrTTL:            30 * time.Second,
	LockTimeTTL:       30 * time.Minute,
	RedisKeyConnector: "exchange:connector",
	RedisKeyLockTime:  "exchange:lockTime",
	RedisKeyErrTimeAt: "exchange:errTime",
}

type Option func(*Config)

func WithErrThreshold(threshold int) Option {
	return func(c *Config) {
		c.ErrThreshold = threshold
	}
}

func WithErrTTL(ttl time.Duration) Option {
	return func(c *Config) {
		c.ErrTTL = ttl
	}
}

func WithLockTimeTTL(ttl time.Duration) Option {
	return func(c *Config) {
		c.LockTimeTTL = ttl
	}
}

func WithRedisKeys(connector, lockTime, errTimeAt string) Option {
	return func(c *Config) {
		c.RedisKeyConnector = connector
		c.RedisKeyLockTime = lockTime
		c.RedisKeyErrTimeAt = errTimeAt
	}
}

type ProxyOption func(*proxyOptions)

type proxyOptions struct {
	primaryConnector   ExchangeConnector
	standbyConnector   ExchangeConnector
	cache              redis.UniversalClient
	alertService       IAlertService
	config             Config
}

func WithPrimaryConnector(c ExchangeConnector) ProxyOption {
	return func(o *proxyOptions) {
		o.primaryConnector = c
	}
}

func WithStandbyConnector(c ExchangeConnector) ProxyOption {
	return func(o *proxyOptions) {
		o.standbyConnector = c
	}
}

func WithCache(c redis.UniversalClient) ProxyOption {
	return func(o *proxyOptions) {
		o.cache = c
	}
}

func WithAlertService(s IAlertService) ProxyOption {
	return func(o *proxyOptions) {
		o.alertService = s
	}
}

func WithConfig(cfg Config) ProxyOption {
	return func(o *proxyOptions) {
		o.config = cfg
	}
}

func NewProxy(opts ...ProxyOption) ExchangeApiProxyImpl {
	options := proxyOptions{
		config: DefaultConfig,
	}
	for _, opt := range opts {
		opt(&options)
	}
	return ExchangeApiProxyImpl{
		BinanceImpl:  options.primaryConnector,
		OKXImpl:      options.standbyConnector,
		Cache:        options.cache,
		AlertService: options.alertService,
	}
}

func NewAdapter(proxy ExchangeApiProxy) ExchangeApi {
	return ExchangeApiAdapter{
		ApiProxy: proxy,
	}
}
