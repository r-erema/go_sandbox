package example1_test

import (
	"testing"
	"time"

	"github.com/r-erema/go_sendbox/learning/software_design_patterns/creational/abstract_factory/example1"
	"github.com/stretchr/testify/assert"
)

func TestFactory(t *testing.T) {
	t.Parallel()

	appLayout := func(factory example1.UIFactory) string {
		datetime := time.Date(2021, time.Month(2), 21, 1, 10, 30, 0, time.UTC)

		return factory.CreateLayout().Render("Clock App", datetime)
	}

	assert.Equal(
		t,
		"<!doctype html><html><head><title>Clock App<title><head><body><h1>2021-02-21 01:10:30 +0000 UTC<h1><footer>⏰</footer><body></html>",
		appLayout(example1.HTMLFactory{}),
	)

	assert.Equal(
		t,
		"==== Clock App ====\n|  2021-02-21 01:10:30 +0000 UTC  |\n===== ⏰ =====",
		appLayout(example1.CliFactory{}),
	)
}
