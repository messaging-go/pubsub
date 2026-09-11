generate:
	go tool mockgen -destination=./test/mocks/mock_closer.go -package=mocks -typed io Closer
	go tool mockgen -destination=./test/mocks/mock_receiver.go -package=mocks -typed github.com/messaging-go/pubsub Receiver
