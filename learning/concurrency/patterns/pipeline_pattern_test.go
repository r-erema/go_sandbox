package patterns_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thoas/go-funk"
)

func complexCalculationsPipeline(amount int) int {
	return <-sum(power(generator(amount)))
}

func generator(maxAmount int) <-chan int {
	out := make(chan int)

	go func() {
		for i := 0; i <= maxAmount; i++ {
			out <- i
		}
		close(out)
	}()

	return out
}

func power(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		for v := range in {
			out <- v * v
		}

		close(out)
	}()

	return out
}

func sum(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		var totalSum int
		for v := range in {
			totalSum += v
		}
		out <- totalSum
		close(out)
	}()

	return out
}

func TestComplexCalculationsPipeline(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		amount, want int
	}{
		{
			name:   "Amount is 3",
			amount: 3,
			want:   14,
		},
		{
			name:   "Amount is 5",
			amount: 5,
			want:   55,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, complexCalculationsPipeline(tt.amount))
		})
	}
}

func machineLearningPipeline(input string) []float64 {
	return <-vectorised(stem(removeStopWords(tokenize(input))))
}

func tokenize(in string) <-chan []string {
	out := make(chan []string)

	go func() {
		out <- strings.Fields(in)
		close(out)
	}()

	return out
}

func removeStopWords(input <-chan []string) <-chan []string {
	out := make(chan []string)

	stopWords := map[string]struct{}{
		"and": {},
		"the": {},
		"is":  {},
		"of":  {},
	}

	go func() {
		for words := range input {
			out <- funk.FilterString(words, func(word string) bool {
				_, ok := stopWords[word]

				return !ok
			})
		}

		close(out)
	}()

	return out
}

func stem(input <-chan []string) <-chan []string {
	stemmingRules := map[string]string{
		"running": "run",
		"flies":   "fly",
	}

	out := make(chan []string)

	go func() {
		for tokens := range input {
			stemmed := make([]string, 0)

			for _, token := range tokens {
				if replacement, found := stemmingRules[token]; found {
					stemmed = append(stemmed, replacement)
				} else {
					stemmed = append(stemmed, token)
				}
			}
			out <- stemmed
		}

		close(out)
	}()

	return out
}

func vectorised(in <-chan []string) <-chan []float64 {
	out := make(chan []float64)
	go func() {
		for token := range in {
			dummyVector := []float64{float64(len(token)) / 10, float64(len(token)) / 11, float64(len(token)) / 12}
			out <- dummyVector
		}

		close(out)
	}()

	return out
}

func TestMachineLearningPipeline(t *testing.T) {
	t.Parallel()

	input := "the quick brown fox jumps over the lazy dog and keeps running"

	assert.Equal(t, []float64{0.9, 0.8181818181818182, 0.75}, machineLearningPipeline(input))
}
