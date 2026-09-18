package pubsub

import (
	"context"

	"cloud.google.com/go/pubsub/v2"
)

type Receiver interface {
	Receive(ctx context.Context, handler func(ctx context.Context, message *pubsub.Message)) error
}

type Handler struct {
	client Receiver
	ctx    context.Context
	cancel context.CancelFunc
}

func (h Handler) Process(
	_ context.Context,
	_ *pubsub.Message,
	next func(ctx context.Context, item *pubsub.Message) error,
) error {
	return h.client.Receive( //nolint:contextcheck,wrapcheck // context is by design and for error,
		// we want to make sure we have small footprint here, so there's no error transformation/wrap for errors
		h.ctx,
		func(ctx context.Context, message *pubsub.Message) {
			err := next(ctx, message)
			if err != nil {
				message.Nack()

				return
			}

			message.Ack()
		})
}

func (h Handler) Wait() {
	<-h.ctx.Done()
}

func (h Handler) Stop() {
	h.cancel()
}

func New(subscriber Receiver) *Handler {
	ctx, cancel := context.WithCancel(context.Background())

	return &Handler{
		client: subscriber,
		ctx:    ctx,
		cancel: cancel,
	}
}
