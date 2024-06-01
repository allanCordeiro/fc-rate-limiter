package cache

import "context"

type Cache interface {
	Get(context.Context, string) (int64, error)
	Incr(context.Context, string) error
	Expire(context.Context, string, int64) error
}
