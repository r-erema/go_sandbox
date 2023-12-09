package context_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	httpEndpointDelayResponse = "/delay-response"
	serviceDelayMilliseconds  = 100
)

func TestContext(t *testing.T) {
	t.Parallel()

	router := http.NewServeMux()
	router.HandleFunc(httpEndpointDelayResponse, func(writer http.ResponseWriter, request *http.Request) {
		body, err := io.ReadAll(request.Body)
		require.NoError(t, err)

		symbols, err := serviceExcludeCommas(request.Context(), body)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)

			return
		}

		writer.WriteHeader(http.StatusOK)
		_, err = writer.Write(symbols)

		require.NoError(t, err)
	})

	server := httptest.NewServer(router)

	tests := []struct {
		name    string
		request func(t *testing.T) *http.Request
		assert  func(t *testing.T, response *http.Response, err error)
	}{
		{
			name: "successful request",
			request: func(t *testing.T) *http.Request {
				t.Helper()
				request, err := http.NewRequestWithContext(
					context.Background(),
					http.MethodPost,
					server.URL+httpEndpointDelayResponse,
					bytes.NewBuffer([]byte("1,2,3,4")),
				)
				require.NoError(t, err)

				return request
			},
			assert: func(t *testing.T, response *http.Response, err error) {
				t.Helper()
				assert.NoError(t, err)
				body, err := io.ReadAll(response.Body)
				require.NoError(t, err)
				assert.Equal(t, []byte{'1', '2', '3', '4'}, body)
			},
		},
		{
			name: "context timeout",
			request: func(t *testing.T) *http.Request {
				t.Helper()
				ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
				t.Cleanup(cancel)
				request, err := http.NewRequestWithContext(
					ctx,
					http.MethodPost,
					server.URL+httpEndpointDelayResponse,
					bytes.NewBuffer([]byte("1,2,3,4,5,6,7,8")),
				)
				require.NoError(t, err)

				return request
			},
			assert: func(t *testing.T, _ *http.Response, err error) {
				t.Helper()
				assert.ErrorIs(t, err, context.DeadlineExceeded)
			},
		},
	}

	for _, tt := range tests {
		testCase := tt
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			request := testCase.request(t)
			t.Cleanup(func() {
				require.NoError(t, request.Body.Close())
			})
			response, err := http.DefaultClient.Do(request)
			t.Cleanup(func() {
				if response != nil {
					require.NoError(t, response.Body.Close())
				}
			})

			testCase.assert(t, response, err)
		})
	}
}

func serviceExcludeCommas(ctx context.Context, input []byte) ([]byte, error) {
	symbolSets := bytes.Split(input, []byte{','})
	result := make([]byte, 0)

	addNumberToSlice := func(slice []byte, number []byte) []byte {
		time.Sleep(time.Millisecond * serviceDelayMilliseconds)

		return append(slice, number...)
	}

	for _, symbolSet := range symbolSets {
		select {
		case <-ctx.Done():
			return result, fmt.Errorf("context is done prematurely: %w", ctx.Err())
		default:
			result = addNumberToSlice(result, symbolSet)
		}
	}

	return result, nil
}
