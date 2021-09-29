package example1

import (
	"strings"
	"time"
)

type uiRenderer interface {
	factory() ui
	RenderUI(currentTime time.Time) string
}

type CommonRenderer struct {
	concreteRenderer uiRenderer
}

func (cr *CommonRenderer) factory() ui {
	return cr.concreteRenderer.factory()
}

func (cr CommonRenderer) RenderUI(t time.Time) string {
	return strings.Replace(cr.factory().template(), "{{ time }}", t.String(), 1)
}

type CliUICreator struct {
	commonRenderer *CommonRenderer
}

func NewCliUICreator() *CliUICreator {
	commonRenderer := &CommonRenderer{concreteRenderer: nil}
	cliCreator := &CliUICreator{commonRenderer: nil}
	commonRenderer.concreteRenderer = cliCreator
	cliCreator.commonRenderer = commonRenderer

	return cliCreator
}

func (c CliUICreator) factory() ui {
	return cliUI{}
}

func (c CliUICreator) RenderUI(t time.Time) string {
	return c.commonRenderer.RenderUI(t)
}

type WebUICreator struct {
	commonRenderer *CommonRenderer
}

func NewWebUICreator() *WebUICreator {
	commonRenderer := &CommonRenderer{concreteRenderer: nil}
	webCreator := &WebUICreator{commonRenderer: nil}
	commonRenderer.concreteRenderer = webCreator
	webCreator.commonRenderer = commonRenderer

	return webCreator
}

func (c WebUICreator) factory() ui {
	return webUI{}
}

func (c WebUICreator) RenderUI(currentTime time.Time) string {
	return c.commonRenderer.RenderUI(currentTime)
}
