package patterns_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type betOrder struct {
	userID,
	matchID,
	betType string
	amount,
	odds,
	potentialWin float64
}

var errInvalidBetAmount = errors.New("invalid bet amount")

type betOrderProcessor struct {
	ordersCh chan *betOrder
	wg       *sync.WaitGroup
}

func (op *betOrderProcessor) receiveOrders(order *betOrder) error {
	if order.amount <= 0 {
		return errInvalidBetAmount
	}

	op.wg.Add(1)
	op.ordersCh <- order

	return nil
}

func (op *betOrderProcessor) processOrders() {
	for o := range op.ordersCh {
		o.potentialWin = o.amount * o.odds

		op.wg.Done()
	}
}

func TestBetting(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		order         *betOrder
		expectedError error
		expectedWin   float64
	}{
		{
			name: "place valid order",
			order: &betOrder{
				userID:  "user123",
				matchID: "match456",
				betType: "WIN",
				amount:  100,
				odds:    1.5,
			},
			expectedError: nil,
			expectedWin:   150,
		},
		{
			name: "place invalid order",
			order: &betOrder{
				userID:  "user987",
				matchID: "match456",
				betType: "WIN",
				amount:  -100,
				odds:    1.5,
			},
			expectedError: errInvalidBetAmount,
			expectedWin:   0,
		},
	}

	processor := betOrderProcessor{
		ordersCh: make(chan *betOrder, 10),
		wg:       new(sync.WaitGroup),
	}
	go processor.processOrders()

	for _, tt := range tests {
		testCase := tt
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			err := processor.receiveOrders(testCase.order)
			processor.wg.Wait()

			assert.ErrorIs(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedWin, testCase.order.potentialWin)
		})
	}
}

type item struct {
	data int
}

type streamer struct {
	ch chan item
	wg *sync.WaitGroup
}

func buildStreamer(bufferSize int) *streamer {
	return &streamer{
		ch: make(chan item, bufferSize),
		wg: new(sync.WaitGroup),
	}
}

func (s *streamer) producer(items []item) {
	s.wg.Add(1)

	go func() {
		for _, itm := range items {
			s.ch <- itm
		}

		s.wg.Done()
		close(s.ch)
	}()
}

func (s *streamer) consumer(processedItems *[]item) {
	s.wg.Add(1)

	go func() {
		for itm := range s.ch {
			itm.data *= 2
			*processedItems = append(*processedItems, itm)
		}

		s.wg.Done()
	}()
}

func (s *streamer) wait() {
	s.wg.Wait()
}

func TestProducerConsumer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		itemsToProduce []item
		expected       []item
	}{
		{
			name:           "5 items",
			itemsToProduce: []item{{1}, {2}, {3}, {4}, {5}},
			expected:       []item{{2}, {4}, {6}, {8}, {10}},
		},
		{
			name:           "3 items",
			itemsToProduce: []item{{6}, {7}, {8}},
			expected:       []item{{12}, {14}, {16}},
		},
	}

	for _, tt := range tests {
		testCase := tt
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			stream := buildStreamer(10)
			processedItems := make([]item, 0)
			stream.producer(testCase.itemsToProduce)
			stream.consumer(&processedItems)
			stream.wait()

			assert.Equal(t, testCase.expected, processedItems)
		})
	}
}
