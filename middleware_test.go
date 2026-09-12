package pubsub_test

import (
	"context"
	"fmt"
	"testing"
	"testing/synctest"
	"time"

	pubsub2 "cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	pubsubContainer "github.com/testcontainers/testcontainers-go/modules/gcloud/pubsub"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/mock/gomock"

	"github.com/messaging-go/core"
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

func TestNew(t *testing.T) {
	t.Run("processes all messages", func(t *testing.T) {
		uri := createPubSubEmulator(t)
		t.Setenv("PUBSUB_EMULATOR_HOST", uri)
		topic, subscription := createTopicAndSubscription(t, "my-topic", "subscription")
		// populate million message first
		client, err := pubsub2.NewClient(t.Context(), projectID)
		require.NoError(t, err)
		const count = 30_000
		publishMessages(t, client, topic, count)
		t.Log("published messages, processing next...")
		// now we'd like to test that it truly will consume all that message
		processor := core.New[pubsub2.Message](func(err error) {
			// do nothing
		})
		mw := pubsub.New(client.Subscriber(subscription))
		processor.AddMiddleware(mw)
		done := make(chan struct{}, count)
		go func() {
			for {
				if len(done) == count {
					mw.Stop()
					processor.Stop()
					time.Sleep(time.Millisecond * 50)
					break
				}
			}
		}()
		startTime := time.Now()
		processor.Run(func(ctx context.Context, item *pubsub2.Message) error {
			assert.Equal(t, []byte("hello world"), item.Data)
			done <- struct{}{}

			return nil
		})
		t.Log("took ", time.Since(startTime))
	})
}

func publishMessages(tb testing.TB, client *pubsub2.Client, topic string, count int) {
	tb.Helper()
	publisher := client.Publisher(topic)
	for range count {
		publisher.Publish(tb.Context(), &pubsub2.Message{
			Data: []byte("hello world"),
		})
	}
	publisher.Flush()
}

func createPubSubEmulator(tb testing.TB) string {
	tb.Helper()
	pubsubExec, err := pubsubContainer.Run(
		tb.Context(),
		"gcr.io/google.com/cloudsdktool/cloud-sdk:582.0.0-emulators",
		pubsubContainer.WithProjectID(projectID),
		testcontainers.WithExposedPorts("8681/tcp"),
		testcontainers.WithWaitStrategyAndDeadline(time.Minute, wait.ForLog("listening on").WithPollInterval(time.Second)),
	)
	require.NoError(tb, err)

	tb.Cleanup(func() {
		require.NoError(tb, pubsubExec.Terminate(context.Background()))
	})

	return pubsubExec.URI()
}

func createTopicAndSubscription(tb testing.TB, topic, subscription string) (string, string) {
	tb.Helper()
	client, err := pubsub2.NewClient(
		tb.Context(),
		projectID,
	)
	require.NoError(tb, err)
	topicRes, err := client.TopicAdminClient.CreateTopic(tb.Context(), &pubsubpb.Topic{Name: fmt.Sprintf("projects/%s/topics/%s", projectID, topic)})
	require.NoError(tb, err)
	subRes, err := client.SubscriptionAdminClient.CreateSubscription(tb.Context(), &pubsubpb.Subscription{
		Name:  fmt.Sprintf("projects/%s/subscriptions/%s", projectID, subscription),
		Topic: topicRes.Name,
	})
	require.NoError(tb, err)

	return topicRes.Name, subRes.Name
}
