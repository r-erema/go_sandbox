package example1

type ui interface {
	template() string
}

type cliUI struct{}

func (c cliUI) template() string {
	return "==== Current Time =====\n|  {{ time }}  |\n============="
}

type webUI struct{}

func (w webUI) template() string {
	return "<html><title>Current Time</title><body><h1>{{ time }}</h1></body></html>"
}
