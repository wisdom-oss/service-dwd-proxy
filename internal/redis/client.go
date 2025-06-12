package redis

import (
	"context"
	"errors"
	"strings"

	lib "github.com/redis/go-redis/v9"

	"microservice/internal"
)

// This package provides a connection to a redis database and some functions for
// reading and writing data into the redis database

var errEmptyUri = errors.New("empty redis uri in config")
var errRedisPingFailed = errors.New("unable to ping redis database")

var client *lib.Client

func Connect() error {
	uri := internal.Configuration().GetString(internal.ConfigKey_RedisURI)
	if strings.TrimSpace(uri) == "" {
		return errEmptyUri
	}

	opts, err := lib.ParseURL(uri)
	if err != nil {
		return err
	}

	client = lib.NewClient(opts)
	if client.Ping(context.Background()).Err() != nil {
		return errRedisPingFailed
	}
	return nil
}

func Client() *lib.Client {
	return client
}

func IsNotFound(err error) bool {
	return errors.Is(err, lib.Nil)
}
