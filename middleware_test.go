package pubsub_test

import (
	"context"
	"testing"
	"testing/synctest"
	"time"

	pubsub2 "cloud.google.com/go/pubsub/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/messaging-go/pubsub"
	"github.com/messaging-go/pubsub/test/mocks"
)

func TestHandler_Process(t *testing.T) {
	t.Parallel()
	t.Run("returns no error when processing", func(t *testing.T) {
		t.Parallel()

		mockReceiver := mocks.NewMockReceiver(gomock.NewController(t))
		mockReceiver.
			EXPECT().
			Receive(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, f func(context.Context, *pubsub2.Message)) error {
				f(ctx, &pubsub2.Message{
					Data: []byte("hello world"),
				})

				return nil
			}).AnyTimes()
		handler := pubsub.New(mockReceiver)

		assert.NoError(t, handler.Process(t.Context(), nil, func(ctx context.Context, item *pubsub2.Message) error {
			assert.Equal(t, "hello world", string(item.Data))

			return nil
		}))
		assert.NoError(t, handler.Process(t.Context(), nil, func(ctx context.Context, item *pubsub2.Message) error {
			assert.Equal(t, "hello world", string(item.Data))

			return assert.AnError
		}))
	})
	t.Run("cancellation is supported", func(t *testing.T) {
		t.Parallel()

		synctest.Test(t, func(t *testing.T) {
			mockReceiver := mocks.NewMockReceiver(gomock.NewController(t))
			handler := pubsub.New(mockReceiver)

			go func() {
				time.Sleep(time.Minute)
				handler.Stop()
			}()

			now := time.Now()

			handler.Wait()
			assert.Equal(t, time.Minute, time.Since(now))
		})
	})
}
