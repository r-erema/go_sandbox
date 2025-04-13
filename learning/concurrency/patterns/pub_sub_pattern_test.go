package patterns_test

import (
	"fmt"
	"log"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStockExchange(t *testing.T) {
	t.Parallel()

	nasdaq := newStockExchange()
	nyse := newStockExchange()

	bot1 := &tradingBot{id: "bot 1"}
	bot2 := &tradingBot{id: "bot 2"}

	t.Run("multiple publishers and subscribers", func(t *testing.T) {
		t.Parallel()

		nasdaq.addSubscriber(bot1)
		nasdaq.addSubscriber(bot2)
		nyse.addSubscriber(bot1)
		nyse.addSubscriber(bot2)

		nasdaq.updatePrice("AAPL", 150.0)
		assert.Equal(t, "AAPL", bot1.lastTicket)
		assert.InEpsilon(t, 150.0, bot1.lastPrice, 0)
		assert.Equal(t, "AAPL", bot2.lastTicket)
		assert.InEpsilon(t, 150.0, bot2.lastPrice, 0)

		nyse.updatePrice("GOOGL", 150.0)
		assert.Equal(t, "GOOGL", bot1.lastTicket)
		assert.InEpsilon(t, 150.0, bot1.lastPrice, 0)
		assert.Equal(t, "GOOGL", bot2.lastTicket)
		assert.InEpsilon(t, 150.0, bot2.lastPrice, 0)

		nasdaq.removeSubscriber(bot1)
		nyse.removeSubscriber(bot2)
	})
}

type stockExchange struct {
	tickers     map[string]float64
	subscribers map[stockSubscriber]struct{}
	mu          sync.Mutex
}

func newStockExchange() *stockExchange {
	return &stockExchange{
		tickers:     make(map[string]float64),
		subscribers: make(map[stockSubscriber]struct{}),
	}
}

func (s *stockExchange) addSubscriber(subscriber stockSubscriber) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.subscribers[subscriber] = struct{}{}
}

func (s *stockExchange) removeSubscriber(subscriber stockSubscriber) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.subscribers, subscriber)
}

func (s *stockExchange) notify(ticker string, price float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for listener := range s.subscribers {
		listener.onPriceUpdate(ticker, price)
	}
}

func (s *stockExchange) updatePrice(ticker string, price float64) {
	s.mu.Lock()
	s.tickers[ticker] = price
	s.mu.Unlock()
	s.notify(ticker, price)
}

type stockSubscriber interface {
	onPriceUpdate(ticker string, price float64)
}

type tradingBot struct {
	id,
	lastTicket string
	lastPrice float64
}

func (b *tradingBot) onPriceUpdate(ticker string, price float64) {
	b.lastTicket = ticker
	b.lastPrice = price
}

type (
	notificationsSubscriber interface {
		processMessage(msg string) error
		close()
	}
)

type messagesPublisher struct {
	subscribers map[notificationsSubscriber]struct{}
	addSubCh    chan notificationsSubscriber
	removeSubCh chan notificationsSubscriber
	inCh        chan string
	stopCh      chan struct{}
}

func newMessagesPublisher() *messagesPublisher {
	return &messagesPublisher{
		subscribers: make(map[notificationsSubscriber]struct{}),
		addSubCh:    make(chan notificationsSubscriber),
		removeSubCh: make(chan notificationsSubscriber),
		inCh:        make(chan string),
		stopCh:      make(chan struct{}),
	}
}

func (mp *messagesPublisher) start() {
	for {
		select {
		case msg := <-mp.inCh:
			for subscriber := range mp.subscribers {
				err := subscriber.processMessage(msg)
				if err != nil {
					log.Printf("subscriber processMessage error: %s", err)
				}
			}
		case sub := <-mp.addSubCh:
			mp.subscribers[sub] = struct{}{}
		case sub := <-mp.removeSubCh:
			if _, ok := mp.subscribers[sub]; ok {
				sub.close()
				delete(mp.subscribers, sub)
			}
		case <-mp.stopCh:
			close(mp.addSubCh)
			close(mp.inCh)
			close(mp.removeSubCh)

			return
		}
	}
}

func (mp *messagesPublisher) addSubscriberCh() chan<- notificationsSubscriber {
	return mp.addSubCh
}

func (mp *messagesPublisher) removeSubscriberCh() chan<- notificationsSubscriber {
	return mp.removeSubCh
}

func (mp *messagesPublisher) publishingCh() chan<- string {
	return mp.inCh
}

func (mp *messagesPublisher) stop() {
	close(mp.stopCh)
}

type messagesSubscriber struct {
	id     string
	result chan string
}

func newMessagesSubscriber(id string) *messagesSubscriber {
	return &messagesSubscriber{id: id, result: make(chan string)}
}

func (ms *messagesSubscriber) processMessage(msg string) error {
	ms.result <- fmt.Sprintf("subscriber `%s`, processed message: %s", ms.id, msg)

	return nil
}

func (ms *messagesSubscriber) close() {
	close(ms.result)
}

func TestPublisher(t *testing.T) {
	t.Parallel()

	pub, sub, sub2 := newMessagesPublisher(), newMessagesSubscriber(
		"sub 1",
	), newMessagesSubscriber(
		"sub 2",
	)

	go pub.start()

	pub.addSubscriberCh() <- sub
	pub.addSubscriberCh() <- sub2

	pub.publishingCh() <- "hello"

	assert.Equal(t, "subscriber `sub 1`, processed message: hello", <-sub.result)
	assert.Equal(t, "subscriber `sub 2`, processed message: hello", <-sub2.result)

	pub.removeSubscriberCh() <- sub
	pub.publishingCh() <- "after remove"

	_, ok := <-sub.result
	assert.False(t, ok)
	assert.Equal(t, "subscriber `sub 2`, processed message: after remove", <-sub2.result)

	pub.stop()
}
