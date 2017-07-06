package logging

import (
	stdcontex "context"

	"github.com/go-kit/kit/log"
)

var DefaultKeys []interface{} = []interface{}{
	"X-Request-Id",
	"X-Client-Id",
	"X-Subject-Id",
}

func WithContext(logger log.Logger, ctx stdcontex.Context) log.Logger {
	return WithContextKeys(logger, ctx, DefaultKeys)
}

func WithContextKeys(logger log.Logger, ctx stdcontex.Context, keys []interface{}) log.Logger {
	keyvals := make([]interface{}, 0, 0)

	for _, key := range keys {
		if val := ctx.Value(key); val != nil {
			keyvals = append(keyvals, key, val)
		}
	}

	return log.With(logger, keyvals...)
}
