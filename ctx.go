package wits

import (
	"context"
)

var appCtx context.Context

func init() {
	appCtx = context.WithValue(context.Background(), "version", version)
}

func SetCtx(ctx context.Context) {
	appCtx = ctx
}
