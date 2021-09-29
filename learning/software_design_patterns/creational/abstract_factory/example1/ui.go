package example1

import (
	"fmt"
	"time"
)

type header interface {
	render(title string) string
}

type body interface {
	render(time time.Time) string
}

type footer interface {
	render() string
}

type Layout interface {
	Render(title string, time time.Time) string
}

type HTMLHeader struct{}

func (h HTMLHeader) render(title string) string {
	return fmt.Sprintf("<head><title>%s<title><head>", title)
}

type HTMLBody struct {
	footer HTMLFooter
}

func (h HTMLBody) render(t time.Time) string {
	return fmt.Sprintf("<body><h1>%s<h1>%s<body>", t.String(), h.footer.render())
}

type HTMLFooter struct{}

func (h HTMLFooter) render() string {
	return "<footer>⏰</footer>"
}

type HTMLLayout struct {
	header HTMLHeader
	body   HTMLBody
}

func (h HTMLLayout) Render(title string, t time.Time) string {
	return fmt.Sprintf("<!doctype html><html>%s%s</html>", h.header.render(title), h.body.render(t))
}

type CliHeader struct{}

func (c CliHeader) render(title string) string {
	return fmt.Sprintf("==== %s ====\n", title)
}

type CliBody struct{}

func (c CliBody) render(t time.Time) string {
	return fmt.Sprintf("|  %s  |\n", t.String())
}

type CliFooter struct{}

func (c CliFooter) render() string {
	return "===== ⏰ ====="
}

type CliLayout struct {
	header
	body
	footer
}

func (c CliLayout) Render(title string, t time.Time) string {
	return fmt.Sprintf("%s%s%s", c.header.render(title), c.body.render(t), c.footer.render())
}
