package globals

import (
	"github.com/redis/go-redis/v9"
)

// This file contains all globally shared connections (e.g., Databases)

// RedisClient contains the redis client used for accessing the discovery
// results stored in redis for caching purposes
var RedisClient *redis.Client
