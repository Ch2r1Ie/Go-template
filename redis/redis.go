package redis

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	defaultReadTimeout = 2 * time.Second
	defaultMaxRetries  = 2
)

func New(addr string, password string) redis.UniversalClient {
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:       []string{addr},
		Password:    password, // no password set
		ReadTimeout: defaultReadTimeout,
		MaxRetries:  defaultMaxRetries,
		DB:          0, // use default DB
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Panic(err)
	}

	return rdb
}

func NewCluster(addrs []string) *redis.ClusterClient {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:         addrs,
		ReadTimeout:   defaultReadTimeout,
		MaxRetries:    defaultMaxRetries,
		RouteRandomly: true,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := rdb.ForEachShard(ctx, func(ctx context.Context, shard *redis.Client) error {
		return shard.Ping(ctx).Err()
	})
	if err != nil {
		log.Panic(err)
	}

	return rdb
}

func NewFailOver(masterName string, addrs []string) *redis.Client { // Sentinel
	rdb := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    masterName,
		SentinelAddrs: addrs,
		ReadTimeout:   defaultReadTimeout,
		MaxRetries:    defaultMaxRetries,
		RouteRandomly: true,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Panic(err)
	}

	return rdb
}
