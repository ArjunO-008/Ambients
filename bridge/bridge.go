package bridge

import "context"

type Bridge struct {
	ctx context.Context
}

func NewBridge() *Bridge {
	return &Bridge{}
}

func (b *Bridge) SetContext(ctx context.Context) {
	b.ctx = ctx
}
