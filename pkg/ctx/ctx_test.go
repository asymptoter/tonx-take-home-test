package ctx

import (
	"testing"
)

func TestCTX_With(t *testing.T) {
	c := Background()
	c.With("key1", "value1", "key2", "value2").Info("OK")
}
