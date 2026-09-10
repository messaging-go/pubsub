package pubsub_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	pubsub "github.com/messaging-go/pubsub"
)

func TestHandler_Process(t *testing.T) {
	t.Parallel()
	t.Run("panics with not implemented message", func(t *testing.T) {
		t.Parallel()

		handler := pubsub.New[int]()

		assert.Panics(t, func() {
			_ = handler.Process(t.Context(), 1, nil) //nolint:errcheck // this would not return
		})
	})
}
