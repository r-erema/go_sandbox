package example1_test

import (
	"testing"
	"time"

	"github.com/r-erema/go_sendbox/learning/software_design_patterns/creational/factory_method/example1"
	"github.com/stretchr/testify/assert"
)

func TestFactoryMethod(t *testing.T) {
	t.Parallel()

	datetime := time.Date(2021, time.Month(2), 21, 1, 10, 30, 0, time.UTC)
	assert.Equal(
		t,
		"==== Current Time =====\n|  2021-02-21 01:10:30 +0000 UTC  |\n=============",
		example1.NewCliUICreator().RenderUI(datetime),
	)
	assert.Equal(
		t,
		"<html><title>Current Time</title><body><h1>2021-02-21 01:10:30 +0000 UTC</h1></body></html>",
		example1.NewWebUICreator().RenderUI(datetime),
	)
}
