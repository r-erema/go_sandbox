package sync_test

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/rand"
)

type user struct {
	hasPrize bool
}

type tt struct {
	name     string
	testFlow func() []user
}

func busyWaittt() tt {
	return tt{
		name: "Bad approach: busy-wait",
		testFlow: func() []user {
			loggedInUsers := make([]user, 0)
			var mutex sync.Mutex

			outerWorld := func() {
				for {
					time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))

					mutex.Lock()
					loggedInUsers = append(loggedInUsers, user{})
					mutex.Unlock()

					if len(loggedInUsers) >= 100 {
						return
					}
				}
			}

			go outerWorld()

			for {
				mutex.Lock()
				if len(loggedInUsers) >= 100 {
					givePrizes(loggedInUsers[:10])

					return loggedInUsers
				}
				mutex.Unlock()
			}
		},
	}
}

func blockByChanneltt() tt {
	return tt{
		name: "Better approach: block by channel",
		testFlow: func() []user {
			loggedInUsers := make([]user, 0)
			var mutex sync.Mutex

			usersReady := make(chan struct{})

			outerWorld := func() {
				for {
					time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))

					mutex.Lock()
					loggedInUsers = append(loggedInUsers, user{})
					mutex.Unlock()

					if len(loggedInUsers) >= 100 {
						usersReady <- struct{}{}

						return
					}
				}
			}

			go outerWorld()

			<-usersReady

			givePrizes(loggedInUsers[:10])

			return loggedInUsers
		},
	}
}

func syncCondtt() tt {
	return tt{
		name: "Even better approach: sync.Cond",
		testFlow: func() []user {
			loggedInUsers := make([]user, 0)

			cond := sync.NewCond(&sync.Mutex{})

			outerWorld := func() {
				for {
					time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))

					cond.L.Lock()
					loggedInUsers = append(loggedInUsers, user{})
					cond.L.Unlock()

					if len(loggedInUsers) >= 100 {
						cond.Broadcast()

						return
					}
				}
			}

			go outerWorld()

			cond.L.Lock()
			for len(loggedInUsers) < 100 {
				cond.Wait()
			}
			givePrizes(loggedInUsers[:10])
			cond.L.Unlock()

			return loggedInUsers
		},
	}
}

func TestPrizeFirst10LoggedInUsers(t *testing.T) {
	t.Parallel()
	rand.Seed(time.Now().UnixNano())

	tests := []struct {
		name     string
		testFlow func() []user
	}{
		busyWaittt(),
		blockByChanneltt(),
		syncCondtt(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			loggedInUsers := tt.testFlow()

			for i := range loggedInUsers {
				if i < 10 {
					assert.True(t, loggedInUsers[i].hasPrize)
				} else {
					assert.False(t, loggedInUsers[i].hasPrize)
				}
			}
		})
	}
}

func TestSyncCondBroadcast(t *testing.T) {
	t.Parallel()

	cond := sync.NewCond(&sync.Mutex{})

	tests := []struct {
		name                    string
		expectedGoroutinesCount int64
	}{
		{
			name:                    "Call cond.Broadcast() for multiple goroutines",
			expectedGoroutinesCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var runningGoroutinesCount int64

			var runningGoroutinesWG, finishedGoroutinesWG sync.WaitGroup

			for range 4 {
				runningGoroutinesWG.Add(1)
				finishedGoroutinesWG.Add(1)

				go func() {
					cond.L.Lock()

					runningGoroutinesCount++

					runningGoroutinesWG.Done()
					cond.Wait()

					runningGoroutinesCount--

					cond.L.Unlock()
					finishedGoroutinesWG.Done()
				}()
			}

			runningGoroutinesWG.Wait()

			cond.L.Lock()
			cond.Broadcast()
			cond.L.Unlock()

			finishedGoroutinesWG.Wait()

			assert.Equal(t, tt.expectedGoroutinesCount, runningGoroutinesCount)
		})
	}
}

func givePrizes(users []user) {
	for i := range users {
		users[i].hasPrize = true
	}
}
