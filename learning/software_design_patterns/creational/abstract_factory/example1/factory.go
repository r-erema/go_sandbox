package example1

type UIFactory interface {
	createHeader() header
	createBody() body
	createFooter() footer
	CreateLayout() Layout
}

type HTMLFactory struct{}

func (h HTMLFactory) createHeader() header {
	return HTMLHeader{}
}

func (h HTMLFactory) createBody() body {
	return HTMLBody{footer: HTMLFooter{}}
}

func (h HTMLFactory) createFooter() footer {
	return HTMLFooter{}
}

func (h HTMLFactory) CreateLayout() Layout {
	return HTMLLayout{
		header: HTMLHeader{},
		body:   HTMLBody{footer: HTMLFooter{}},
	}
}

type CliFactory struct{}

func (h CliFactory) createHeader() header {
	return CliHeader{}
}

func (h CliFactory) createBody() body {
	return CliBody{}
}

func (h CliFactory) createFooter() footer {
	return CliFooter{}
}

func (h CliFactory) CreateLayout() Layout {
	return CliLayout{
		header: CliHeader{},
		body:   CliBody{},
		footer: CliFooter{},
	}
}
