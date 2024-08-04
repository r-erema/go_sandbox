package example1_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/r-erema/go_sendbox/learning/security/example1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	userID         = "4caf89bf-6122-4770-8fea-2269c55b1d59"
	userPassword   = "secret_password"
	bankAccountPIN = "1427"

	authServerAddr     = "localhost:6606"
	resourceServerAddr = "localhost:6607"
)

type tablett struct {
	name           string
	expectedStatus int
	expectedBody   string
	requestFactory func() *http.Request
}

func initTestStorage() map[example1.UserID]example1.BankAccountPIN {
	storage := make(map[example1.UserID]example1.BankAccountPIN, 1)
	storage[userID] = bankAccountPIN

	return storage
}

func Test_ClientCredentialsGrantFlow(t *testing.T) {
	t.Parallel()

	tests := []tablett{
		ttWithoutCredentials(t),
		ttWithBadCredentials(t),
		ttSuccessful(t),
	}

	authServerStorage := store.NewClientStore()
	err := authServerStorage.Set(userID, &models.Client{
		ID:     userID,
		Secret: userPassword,
		Domain: "",
		UserID: "",
	})
	require.NoError(t, err)

	resourceServer := example1.NewResourceServer(initTestStorage(), fmt.Sprintf("http://%s/interception-endpoint", authServerAddr))

	oAuth2Manager := manage.NewDefaultManager()
	oAuth2Manager.MapClientStorage(authServerStorage)
	oAuth2Manager.MustTokenStorage(store.NewMemoryTokenStore())

	oAuth2 := server.NewDefaultServer(oAuth2Manager)
	oAuth2.SetAllowGetAccessRequest(true)
	oAuth2.SetClientInfoHandler(server.ClientFormHandler)

	authServer := example1.NewAuthorizationServer(oAuth2)

	go func() {
		err := resourceServer.Run(resourceServerAddr)
		assert.NoError(t, err)
	}()

	go func() {
		err := authServer.Run(authServerAddr)
		assert.NoError(t, err)
	}()

	time.Sleep(time.Millisecond * 100)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := tt.requestFactory()

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)

			defer func() {
				err = resp.Body.Close()
				require.NoError(t, err)
			}()

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			assert.Equal(t, tt.expectedBody, strings.TrimSpace(string(body)))
		})
	}
}

func ttWithoutCredentials(t *testing.T) tablett {
	t.Helper()

	return tablett{
		name:           "request without credentials",
		expectedStatus: http.StatusUnauthorized,
		expectedBody:   "",
		requestFactory: func() *http.Request {
			request, err := http.NewRequestWithContext(
				context.Background(),
				http.MethodGet, fmt.Sprintf("http://%s/bank-account-pin", resourceServerAddr),
				http.NoBody,
			)
			require.NoError(t, err)

			return request
		},
	}
}

func ttWithBadCredentials(t *testing.T) tablett {
	t.Helper()

	return tablett{
		name:           "request with bad credentials",
		expectedStatus: http.StatusUnauthorized,
		expectedBody:   "",
		requestFactory: func() *http.Request {
			var err error

			tokenRequest, err := http.NewRequestWithContext(context.Background(),
				http.MethodGet,
				fmt.Sprintf(
					"http://%s/token?grant_type=client_credentials&client_id=%s&client_secret=%s&scope=all",
					authServerAddr,
					userID,
					"bad_password",
				),
				http.NoBody,
			)
			require.NoError(t, err)

			tokenResponse, err := http.DefaultClient.Do(tokenRequest)
			require.NoError(t, err)

			defer func() {
				err = tokenResponse.Body.Close()
				require.NoError(t, err)
			}()

			var data map[string]string
			err = json.NewDecoder(tokenResponse.Body).Decode(&data)

			request, err := http.NewRequestWithContext(
				context.Background(),
				http.MethodGet,
				fmt.Sprintf("http://%s/bank-account-pin?uid=%s", resourceServerAddr, userID),
				http.NoBody,
			)
			require.NoError(t, err)

			request.Header.Add("Authorization", "Bearer "+data["access_token"])

			return request
		},
	}
}

func ttSuccessful(t *testing.T) tablett {
	t.Helper()

	return tablett{
		name:           "successful request",
		expectedStatus: http.StatusOK,
		expectedBody:   bankAccountPIN,
		requestFactory: func() *http.Request {
			tokenRequest, err := http.NewRequestWithContext(context.Background(),
				http.MethodGet,
				fmt.Sprintf(
					"http://%s/token?grant_type=client_credentials&client_id=%s&client_secret=%s&scope=all",
					authServerAddr,
					userID,
					userPassword,
				),
				http.NoBody,
			)
			require.NoError(t, err)

			tokenResponse, err := http.DefaultClient.Do(tokenRequest)
			require.NoError(t, err)

			defer func() {
				err = tokenResponse.Body.Close()
				require.NoError(t, err)
			}()

			var data map[string]string
			err = json.NewDecoder(tokenResponse.Body).Decode(&data)

			request, err := http.NewRequestWithContext(
				context.Background(),
				http.MethodGet,
				fmt.Sprintf("http://%s/bank-account-pin?uid=%s", resourceServerAddr, userID),
				http.NoBody,
			)
			require.NoError(t, err)

			request.Header.Add("Authorization", "Bearer "+data["access_token"])

			return request
		},
	}
}
