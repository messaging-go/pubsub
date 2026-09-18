# pubsub

A [messaging-go/core](https://github.com/messaging-go/core) middleware that adds support for consuming messages
from [Google Cloud Pub/Sub](https://cloud.google.com/pubsub) subscriptions.

## Installation

```shell
go get github.com/messaging-go/pubsub
```

## Usage

`pubsub.New` wraps a Pub/Sub subscriber client and returns a `core` middleware. Register it on your
`core.Processor` and call `Run` with your handler — messages are automatically acknowledged when the handler
returns `nil`, and nacked otherwise.

```go
package main

import (
	"context"
	"fmt"
	"log"

	gpubsub "cloud.google.com/go/pubsub/v2"

	"github.com/messaging-go/core"
	"github.com/messaging-go/pubsub"
)

func main() {
	ctx := context.Background()

	client, err := gpubsub.NewClient(ctx, "my-project")
	if err != nil {
		log.Fatal(err)
	}

	subscriber := client.Subscriber("my-subscription")
	// customize subscriber.ReceiveSettings here if needed

	processor := core.New[gpubsub.Message](func(err error) {
		// handle global errors here
		log.Print(err)
	})

	consumerMiddleware := pubsub.New(subscriber)
	processor.AddMiddleware(consumerMiddleware)

	processor.Run(func(ctx context.Context, item *gpubsub.Message) error {
		fmt.Println(string(item.Data))

		return nil
	})
}
```

### Stopping

`Handler` exposes `Stop()` to cancel message consumption and `Wait()` to block until it has stopped, which is
useful for coordinating a graceful shutdown alongside `processor.Stop()`:

```go
go func() {
	<-shutdownSignal
	consumerMiddleware.Stop()
	processor.Stop()
}()
```

See [middleware_example_test.go](middleware_example_test.go) for a full runnable example against the Pub/Sub
emulator.

## How it works

`pubsub.New` accepts anything satisfying the `Receiver` interface (a `*pubsub.Subscriber` from
`cloud.google.com/go/pubsub/v2` satisfies it out of the box):

```go
type Receiver interface {
	Receive(ctx context.Context, handler func(ctx context.Context, message *pubsub.Message)) error
}
```

The returned `Handler` implements `core`'s middleware `Process` method: it calls `Receive` on the underlying
client, invokes the next handler in the chain for each message, and acks/nacks the message based on whether an
error was returned.
