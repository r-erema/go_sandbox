package patterns_test

import (
	"fmt"
	"io"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/util/rand"
)

type poolTask struct {
	id      int
	execute func() error
}

func poolWorker(id int, tasks <-chan poolTask, wg *sync.WaitGroup, stdout io.Writer) {
	for tsk := range tasks {
		err := tsk.execute()
		if err != nil {
			_, _ = fmt.Fprintf(stdout, "error in the task execution in worker %d: %v", id, err)
		}

		wg.Done()
	}
}

func scheduleTasks(workerCount, taskCount int, stdout io.Writer) {
	wg := new(sync.WaitGroup)
	tasks := make(chan poolTask, 10)

	for i := 0; i < workerCount; i++ {
		go poolWorker(i, tasks, wg, stdout)
	}

	for i := 0; i < taskCount; i++ {
		wg.Add(1)
		tasks <- func(i int) poolTask {
			return poolTask{
				id: i,
				execute: func() error {
					time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
					_, _ = fmt.Fprintf(stdout, "Task `%d` is executed\n", i)

					return nil
				},
			}
		}(i)
	}

	wg.Wait()

	close(tasks)
}

func TestTaskManager(t *testing.T) {
	t.Parallel()

	read, write, _ := os.Pipe()

	const tasksCount = 5

	scheduleTasks(3, tasksCount, write)

	err := write.Close()
	require.NoError(t, err)

	out, err := io.ReadAll(read)
	require.NoError(t, err)

	for i := 0; i < tasksCount; i++ {
		assert.Contains(t, string(out), fmt.Sprintf("Task `%d` is executed", i))
	}
}

type gameServerTask struct {
	playerID, action string
}

type workerPool struct {
	tasks           chan gameServerTask
	workerCount     int
	quitChan        chan bool
	tasksInProgress uint8
	wg              *sync.WaitGroup
	cond            *sync.Cond
	stdOut          io.Writer
}

func newWorkerPool(workerCount int, stdOut io.Writer) *workerPool {
	return &workerPool{
		tasks:       make(chan gameServerTask),
		workerCount: workerCount,
		quitChan:    make(chan bool),
		wg:          new(sync.WaitGroup),
		cond:        sync.NewCond(&sync.Mutex{}),
		stdOut:      stdOut,
	}
}

func (wp *workerPool) Start() {
	for i := 0; i < wp.workerCount; i++ {
		wp.wg.Add(1)

		go func() {
			defer wp.wg.Done()

			for {
				select {
				case gameTask := <-wp.tasks:
					wp.cond.L.Lock()
					wp.tasksInProgress++
					wp.cond.L.Unlock()

					processTask(gameTask, wp.stdOut)

					wp.cond.L.Lock()
					wp.tasksInProgress--
					wp.cond.Broadcast()
					wp.cond.L.Unlock()

				case <-wp.quitChan:
					return
				}
			}
		}()
	}

	wp.wg.Wait()
}

func processTask(task gameServerTask, stdout io.Writer) {
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))

	_, _ = fmt.Fprintf(stdout, "Task `%s` for player `%s` is executed\n", task.action, task.playerID)
}

func (wp *workerPool) Stop() {
	defer wp.cond.L.Unlock()

	wp.cond.L.Lock()
	for wp.tasksInProgress > 0 {
		wp.cond.Wait()
	}

	close(wp.quitChan)
}

func TestGameServer(t *testing.T) {
	t.Parallel()

	read, write, _ := os.Pipe()

	pool := newWorkerPool(5, write)
	go pool.Start()

	tasks := []gameServerTask{
		{playerID: "player1", action: "move"},
		{playerID: "player2", action: "jump"},
		{playerID: "player3", action: "shoot"},
	}

	for _, serverTask := range tasks {
		pool.tasks <- serverTask
	}

	pool.Stop()

	err := write.Close()
	require.NoError(t, err)

	out, err := io.ReadAll(read)
	require.NoError(t, err)

	for _, serverTask := range tasks {
		assert.Contains(
			t,
			string(out),
			fmt.Sprintf("Task `%s` for player `%s` is executed\n", serverTask.action, serverTask.playerID),
		)
	}
}
