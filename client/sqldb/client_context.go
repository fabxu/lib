package sqldb

import (
	"context"

	"gorm.io/gorm"
)

type ctxBaseSQL struct{}

var ctxBaseSQLKey = &ctxBaseSQL{}

func (c *Client) Inject(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, ctxBaseSQLKey, tx)
}

func (c *Client) Extract(ctx context.Context, defaultTx *gorm.DB) *gorm.DB {
	tx, ok := ctx.Value(ctxBaseSQLKey).(*gorm.DB)
	if !ok || tx == nil {
		return defaultTx
	}

	return tx
}
