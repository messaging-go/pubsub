package integration_template_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	integration_template "github.com/messaging-go/integration-template"
)

func TestHandler_Process(t *testing.T) {
	t.Parallel()
	t.Run("panics with not implemented message", func(t *testing.T) {
		t.Parallel()

		handler := integration_template.New[int]()

		assert.Panics(t, func() {
			_ = handler.Process(t.Context(), 1, nil) //nolint:errcheck // this would not return
		})
	})
}
