package pubsub_test

import (
	"context"
	"fmt"
	"log"
	"time"

	gpubsub "cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	"github.com/testcontainers/testcontainers-go"
	pubsubContainer "github.com/testcontainers/testcontainers-go/modules/gcloud/pubsub"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/messaging-go/core"
	"github.com/messaging-go/pubsub"
)

const projectID = "messaging-go"

func ExampleNew() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sub, host := prepareEnv(ctx)
	client := mustGetClient(ctx, host)
	subscriptionClient := client.Subscriber(sub.Name)
	// subscriptionClient.ReceiveSettings define your own
	processor := core.New[gpubsub.Message](func(err error) {
		// handle global errors here
		log.Print(err)
	})
	consumerMiddleware := pubsub.New(subscriptionClient)

	go func() {
		time.Sleep(time.Second * 10)
		consumerMiddleware.Stop()
		processor.Stop()
	}()

	processor.AddMiddleware(consumerMiddleware)
	processor.Run(func(ctx context.Context, item *gpubsub.Message) error {
		fmt.Println(string(item.Data))

		return nil
	})
	//Output: Hello messaging-go!
	// Hello messaging-go!
	// Hello messaging-go!
	// Hello messaging-go!
	// Hello messaging-go!
}

func prepareEnv(ctx context.Context) (*pubsubpb.Subscription, string) {
	emulatorHost := preparePubSubInfra(ctx)
	topic, subscription := prepareTopicAndSubscription(ctx, emulatorHost)
	// create
	publishMessageToTopic(ctx, topic, emulatorHost)

	return subscription, emulatorHost
}

func publishMessageToTopic(ctx context.Context, topic *pubsubpb.Topic, host string) {
	client := mustGetClient(ctx, host)

	publisher := client.Publisher(topic.Name)
	for range 5 {
		publisher.Publish(ctx, &gpubsub.Message{
			Data: []byte("Hello messaging-go!"),
		})
	}

	publisher.Flush()
}

func mustGetClient(ctx context.Context, host string) *gpubsub.Client {
	client, err := gpubsub.NewClient(
		ctx,
		projectID,
		option.WithEndpoint(host),
		option.WithoutAuthentication(),
		option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
	)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func prepareTopicAndSubscription(ctx context.Context, host string) (*pubsubpb.Topic, *pubsubpb.Subscription) {
	client := mustGetClient(ctx, host)

	topic, err := client.TopicAdminClient.CreateTopic(ctx, &pubsubpb.Topic{
		Name: fmt.Sprintf("projects/%s/topics/my-topic", projectID),
	})
	if err != nil {
		log.Fatal(err)
	}

	subscription, err := client.SubscriptionAdminClient.CreateSubscription(ctx, &pubsubpb.Subscription{
		Name:  fmt.Sprintf("projects/%s/subscriptions/my-subscription", projectID),
		Topic: topic.Name,
	})
	if err != nil {
		log.Fatal(err)
	}

	return topic, subscription
}

func preparePubSubInfra(ctx context.Context) string {
	pubsubExec, err := pubsubContainer.Run(
		ctx,
		"gcr.io/google.com/cloudsdktool/cloud-sdk:582.0.0-emulators",
		pubsubContainer.WithProjectID(projectID),
		testcontainers.WithExposedPorts("8681/tcp"),
		testcontainers.WithWaitStrategyAndDeadline(time.Minute, wait.ForLog("listening on").WithPollInterval(time.Second)),
	)
	if err != nil {
		log.Fatal(err)
	}

	go func() { //nolint:gosec // this is cleanup func
		<-ctx.Done()
		log.Print("stopping containers")

		cancellationCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		err := pubsubExec.Terminate(cancellationCtx) //nolint:contextcheck // new ctx because the original is done
		if err != nil {
			log.Fatal(err)
		}
	}()

	return pubsubExec.URI()
}
