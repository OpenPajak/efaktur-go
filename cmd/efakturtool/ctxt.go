package main

import (
	"context"

	"github.com/OpenPajak/efaktur-go/pkg/provider/web"
)

type ctxtKey[T any] struct{ _ T }

type (
	ctxtKeyConfT   [0]struct{}
	ctxtKeyClientT [0]struct{}
)

var (
	ctxtKeyConf   ctxtKey[ctxtKeyConfT]
	ctxtKeyClient ctxtKey[ctxtKeyClientT]
)

func (c ctxtKey[T]) Pack(ctx context.Context, data any) context.Context {
	return context.WithValue(ctx, c, data)
}
func (c ctxtKey[T]) Unpack(ctx context.Context) any {
	return ctx.Value(c)
}

func WithConf(ctx context.Context, conf *Config) context.Context {
	return ctxtKeyConf.Pack(ctx, conf)
}
func GetConfFromContext(ctx context.Context) *Config {
	value, ok := ctxtKeyConf.Unpack(ctx).(*Config)
	if !ok || value == nil {
		return nil
	}
	return value
}

func WithClient(ctx context.Context, client *web.Client) context.Context {
	return ctxtKeyClient.Pack(ctx, client)
}
func GetClientFromContext(ctx context.Context) *web.Client {
	value, ok := ctxtKeyClient.Unpack(ctx).(*web.Client)
	if !ok || value == nil {
		return nil
	}
	return value
}
