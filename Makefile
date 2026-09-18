generate:
	go tool mockgen -destination=./test/mocks/mock_receiver.go -package=mocks -typed github.com/messaging-go/pubsub Receiver
