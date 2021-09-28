package example1

import "fmt"

type Item struct {
	observers []observer
	name      string
	inStock   bool
}

func NewItem(name string) *Item {
	return &Item{observers: nil, name: name, inStock: false}
}

func (i *Item) SetAsAvailable() {
	i.inStock = true
	i.notifyAll(fmt.Sprintf("%s is available", i.name))
}

func (i *Item) Register(observer observer) {
	i.observers = append(i.observers, observer)
}

func (i *Item) Deregister(observer observer) {
	length := len(i.observers)

	for j, o := range i.observers {
		if o.ID() == observer.ID() {
			i.observers[length-1], i.observers[j] = i.observers[j], i.observers[length-1]
			i.observers = i.observers[:length-1]

			return
		}
	}
}

func (i *Item) notifyAll(event string) {
	for _, o := range i.observers {
		o.handleEvent(event)
	}
}
