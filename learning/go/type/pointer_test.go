package type_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPointerToFunction(t *testing.T) {
	t.Parallel()

	var ptr *string

	str := "abc123"
	ptr = &str

	assert.Equal(
		t,
		func(ptr *string) string { return fmt.Sprintf("%p", ptr) }(ptr),
		fmt.Sprintf("%p", ptr),
	)
	assert.NotEqual(
		t,
		func(ptr *string) string { return fmt.Sprintf("%p", &ptr) }(ptr),
		fmt.Sprintf("%p", &ptr),
	)
}
