package reltest

import (
	"context"
	"fmt"
)

type ctxKeyType uint8

type ctxData struct {
	txDepth int
}

func (cd ctxData) String() string {
	return fmt.Sprintf("<Transaction: %d>", cd.txDepth)
}

var (
	ctxKey ctxKeyType = 0
)

func fetchContext(ctx context.Context) ctxData {
	if tx, ok := ctx.Value(ctxKey).(ctxData); ok {
		return tx
	}

	return ctxData{}
}

func wrapContext(ctx context.Context, ctxData ctxData) context.Context {
	return context.WithValue(ctx, ctxKey, ctxData)
}
