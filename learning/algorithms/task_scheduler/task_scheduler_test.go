package taskscheduler_test

import (
	"container/heap"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLeastInterval(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		tasks []byte
		n     int
		want  int
	}{
		{
			name:  "1 task",
			tasks: []byte{'A'},
			n:     22,
			want:  1,
		},
		{
			name:  "2 tasks with gap",
			tasks: []byte{'A', 'B', 'B'},
			n:     2,
			want:  4,
		},
		{
			name:  "many tasks",
			tasks: []byte{'A', 'A', 'A', 'B', 'B', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K'},
			n:     7,
			want:  18,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, leastInterval(tt.tasks, tt.n))
		})
	}
}

type tasksCountArr []int

func (t *tasksCountArr) Len() int {
	return len(*t)
}

func (t *tasksCountArr) Less(i, j int) bool {
	return (*t)[i] > (*t)[j]
}

func (t *tasksCountArr) Swap(i, j int) {
	(*t)[i], (*t)[j] = (*t)[j], (*t)[i]
}

func (t *tasksCountArr) Push(x any) {
	taskCount, _ := x.(int)
	*t = append(*t, taskCount)
}

func (t *tasksCountArr) Pop() any {
	old := *t
	n := len(old)
	popped := old[n-1]
	*t = old[0 : n-1]

	return popped
}

// Time O(n), where n is the number of tasks
// Space O(1), since we have at most 26 different characters

func leastInterval(taskNames []byte, n int) int {
	frequencyTmp := make(map[byte]int)
	for i := range taskNames {
		frequencyTmp[taskNames[i]]++
	}

	tasks := make(tasksCountArr, 0, len(frequencyTmp))

	for _, freq := range frequencyTmp {
		tasks = append(tasks, freq)
	}

	heap.Init(&tasks)

	var (
		taskCount, _       = heap.Pop(&tasks).(int)
		queue              = []int{taskCount}
		res, nextTaskCount int
	)

	for ; len(queue) > 0; res++ {
		taskCount, queue = queue[0], queue[1:]
		taskCount--

		for taskCount > 0 && len(queue) < n {
			nextTaskCount = 1
			if len(tasks) > 0 {
				nextTaskCount, _ = heap.Pop(&tasks).(int)
			}

			queue = append(queue, nextTaskCount)
		}

		if taskCount > 0 {
			heap.Push(&tasks, taskCount)
		}

		if len(tasks) > 0 {
			taskCount, _ = heap.Pop(&tasks).(int)
			queue = append(queue, taskCount)
		}
	}

	return res
}
