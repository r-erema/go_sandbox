package patterns_test

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCrawler(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			_, err := writer.Write([]byte(`{"title": "root"}`))
			require.NoError(t, err)
		case "/path-0":
			_, err := writer.Write([]byte(`{"title": "path-0"}`))
			require.NoError(t, err)
		}
	}))

	fch := &httpFetcher{}

	urls := []string{
		server.URL,
		fmt.Sprintf("%s/%s", server.URL, "path-0"),
	}

	results, errors := crawl(urls, fch)
	for _, err := range errors {
		assert.NoError(t, err)
	}

	for _, res := range results {
		assert.Contains(t, res, `"title"`)
	}
}

type fetcher interface {
	fetch(url string) (string, error)
}

type httpFetcher struct{}

func (f *httpFetcher) fetch(url string) (string, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
	if err != nil {
		return "", fmt.Errorf("creation new request error: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP GET error: %w", err)
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Printf("body close error: %s", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading body error: %w", err)
	}

	return string(body), nil
}

func crawl(urls []string, fetcher fetcher) ([]string, []error) {
	var wg sync.WaitGroup

	results, errs := make(chan string, len(urls)), make(chan error, len(urls))

	for _, url := range urls {
		wg.Add(1)

		go func(url string) {
			defer wg.Done()

			body, err := fetcher.fetch(url)
			if err != nil {
				errs <- err

				return
			}
			results <- body
		}(url)
	}

	go func() {
		wg.Wait()
		close(results)
		close(errs)
	}()

	var (
		bodies []string
		errors []error
	)

	for {
		select {
		case body, ok := <-results:
			if ok {
				bodies = append(bodies, body)
			} else {
				results = nil
			}
		case err, ok := <-errs:
			if ok {
				errors = append(errors, err)
			} else {
				errs = nil
			}
		}

		if results == nil && errs == nil {
			break
		}
	}

	return bodies, errors
}

func TestVideoProcessor(t *testing.T) {
	t.Parallel()

	t.Run("process single segment", func(t *testing.T) {
		t.Parallel()

		segment := videoSegment{id: 1}
		encoded := processSegment(segment)
		assert.Equal(t, segment.id, encoded.id)
	})

	t.Run("process multiple segments", func(t *testing.T) {
		t.Parallel()

		segments := []videoSegment{
			{id: 1, frames: []frame{{'1', '8', 'c'}}},
			{id: 2, frames: []frame{{'7'}}},
			{id: 3, frames: []frame{{'o', 'o', '3', 'p', '!'}}},
			{id: 4, frames: []frame{{'r', 'i'}}},
		}
		numWorkers := 2

		encodedSegments := fanOutFanIn(segments, numWorkers)
		assert.Equal(t, len(segments), len(encodedSegments))

		expectedEncodedResults := map[int]string{
			segments[0].id: "43663f6ec7d1d7292c6d4c38545c834d7cf1769745d1ebabacac687fd0c9584d",
			segments[1].id: "7902699be42c8a8e46fbbb4501726517e86b22c56a189f7625a6da49081b2451",
			segments[2].id: "148ca14f9b12d8289968d59b8f37a52e5476d06bb4787db3c82ee55ace3a96b6",
			segments[3].id: "396a14ab206e2b44e03c4e00393e948cce36a6b0f0d7489cb46d944b33ad51c8",
		}
		for _, segment := range encodedSegments {
			assert.Equal(t, expectedEncodedResults[segment.id], fmt.Sprintf("%x", segment.encodedData))
		}
	})
}

type videoSegment struct {
	id     int
	frames []frame
}

type frame []byte

type encodedSegment struct {
	id          int
	encodedData []byte
}

func processSegment(segment videoSegment) encodedSegment {
	var data []byte
	for _, frm := range segment.frames {
		data = append(data, frm...)
	}

	hsh := sha256.New()
	hsh.Write(data)
	encoded := hsh.Sum(nil)

	return encodedSegment{
		id:          segment.id,
		encodedData: encoded,
	}
}

func fanOutFanIn(segments []videoSegment, numWorkers int) []encodedSegment {
	segmentChan := make(chan videoSegment, len(segments))
	encodedChan := make(chan encodedSegment, len(segments))

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for segment := range segmentChan {
				encodedChan <- processSegment(segment)
			}
		}()
	}

	for _, segment := range segments {
		segmentChan <- segment
	}

	close(segmentChan)

	go func() {
		wg.Wait()
		close(encodedChan)
	}()

	encodedSegments := make([]encodedSegment, 0)
	for encoded := range encodedChan {
		encodedSegments = append(encodedSegments, encoded)
	}

	return encodedSegments
}
