package redis

import "time"

type Cache interface {
	GetString(key string) (string, error)
	SetJSON(key string, value interface{}, ttl time.Duration) error
}
