package patterns_test

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type job struct {
	delay     int
	processed bool
}

func (j *job) process() {
	time.Sleep(time.Millisecond * 100)

	j.processed = true
}

func worker(jobs <-chan *job, done <-chan struct{}, wg *sync.WaitGroup) {
	for {
		select {
		case j, ok := <-jobs:
			if ok {
				j.process()
			} else {
				wg.Done()

				return
			}
		case <-done:
			wg.Done()

			return
		}
	}
}

func TestJobProcessing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		numWorkers, numJobs int
	}{
		{name: "1 worker, 5 jobs", numWorkers: 1, numJobs: 5},
		{name: "2 worker, 5 jobs", numWorkers: 2, numJobs: 5},
		{name: "3 worker, 5 jobs", numWorkers: 3, numJobs: 5},
	}

	for _, tt := range tests {
		testCase := tt

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			jobsCh := make(chan *job)
			done := make(chan struct{})

			wg := new(sync.WaitGroup)

			for i := 0; i < testCase.numWorkers; i++ {
				wg.Add(1)
				go worker(jobsCh, done, wg)
			}

			jobs := make([]*job, testCase.numJobs)
			go func() {
				for i := 0; i < testCase.numJobs; i++ {
					jobs[i] = &job{delay: i}
					jobsCh <- jobs[i]
				}
				close(jobsCh)
				done <- struct{}{}
			}()

			wg.Wait()

			for i := range jobs {
				assert.True(t, jobs[i].processed)
			}
		})
	}
}

type task struct {
	id, processedWorkerID int
	processed             bool
}

func taskProducer(tasks chan<- task, done chan<- struct{}, tasksCount int) {
	for i := 0; i < tasksCount; i++ {
		tasks <- task{id: i}
	}
	close(tasks)
	done <- struct{}{}
}

func taskWorker(workerID int, tasks <-chan task, done <-chan struct{}, processedTasks chan<- task, wgWorkers, wgTasks *sync.WaitGroup) {
	for {
		select {
		case tsk, ok := <-tasks:
			if ok {
				time.Sleep(time.Millisecond * 100)

				tsk.processedWorkerID = workerID
				tsk.processed = true
				processedTasks <- tsk

				wgTasks.Done()
			} else {
				wgWorkers.Done()

				return
			}
		case <-done:
			wgWorkers.Done()

			return
		}
	}
}

func TestTaskSystem(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                     string
		tasksCount, workersCount int
	}{
		{name: "2 workers, 6 tasks", tasksCount: 6, workersCount: 2},
		{name: "3 workers, 7 tasks", tasksCount: 7, workersCount: 3},
		{name: "4 workers, 11", tasksCount: 11, workersCount: 4},
	}

	for _, tt := range tests {
		testCase := tt

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			tasks := make(chan task)
			done := make(chan struct{})
			processedTasks := make(chan task, testCase.tasksCount)

			wgWorkers := new(sync.WaitGroup)
			wgTasks := new(sync.WaitGroup)
			wgTasks.Add(testCase.tasksCount)

			go taskProducer(tasks, done, testCase.tasksCount)

			for i := 0; i < testCase.workersCount; i++ {
				wgWorkers.Add(1)
				go taskWorker(i, tasks, done, processedTasks, wgWorkers, wgTasks)
			}

			<-done
			wgWorkers.Wait()

			go func() {
				wgTasks.Wait()
				close(processedTasks)
			}()

			processedTasksCount := 0
			for processedTask := range processedTasks {
				assert.True(t, processedTask.processed)
				processedTasksCount++
			}
			assert.Equal(t, testCase.tasksCount, processedTasksCount)
		})
	}
}
