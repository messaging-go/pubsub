generate:
	go tool mockgen -destination=./test/mocks/mock_closer.go -package=mocks -typed io Closer
