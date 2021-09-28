package example1

import "fmt"

type Customer struct {
	id, name string
	log      []string
}

func NewCustomer(id, name string) *Customer {
	return &Customer{id: id, name: name, log: nil}
}

func (c Customer) ID() string {
	return c.id
}

func (c *Customer) handleEvent(event string) {
	c.log = append(c.log, fmt.Sprintf("I got your message about `%s`, thanks! Best regards, %s", event, c.name))
}

func (c Customer) Log() []string {
	return c.log
}
