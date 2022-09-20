package example1

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type (
	UserID         string
	BankAccountPIN string
)

type ResourceServer struct {
	storage                        map[UserID]BankAccountPIN
	authServerInterceptionEndpoint string
}

func NewResourceServer(storage map[UserID]BankAccountPIN, authServerInterceptionEndpoint string) *ResourceServer {
	return &ResourceServer{storage: storage, authServerInterceptionEndpoint: authServerInterceptionEndpoint}
}

func (rs ResourceServer) Run(addr string) error {
	router := http.NewServeMux()

	router.Handle("/bank-account-pin", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			http.Error(writer, "", http.StatusMethodNotAllowed)

			return
		}

		authResponse, err := rs.authRequest(request.Context(), request)
		if err != nil {
			http.Error(writer, "", http.StatusBadRequest)

			return
		}

		defer func() {
			err = authResponse.Body.Close()
			log.Printf("body close error: %s", err)
		}()

		if authResponse.StatusCode != http.StatusOK {
			http.Error(writer, "", authResponse.StatusCode)

			return
		}

		uid := request.URL.Query().Get("uid")
		pin, ok := rs.storage[UserID(uid)]
		if !ok {
			http.Error(writer, "", http.StatusBadRequest)

			return
		}

		if _, err := writer.Write([]byte(pin)); err != nil {
			http.Error(writer, "", http.StatusInternalServerError)

			return
		}
	}))

	if err := (&http.Server{ //nolint:exhaustruct
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: time.Second,
	}).ListenAndServe(); err != nil {
		return fmt.Errorf("resource server listen and serve error: %w", err)
	}

	return nil
}

func (rs ResourceServer) authRequest(ctx context.Context, request *http.Request) (*http.Response, error) {
	tokenCheckRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		rs.authServerInterceptionEndpoint,
		http.NoBody,
	)
	if err != nil {
		return nil, fmt.Errorf("creation auth tokenCheckRequest error: %w", err)
	}

	tokenCheckRequest.Header = request.Header

	response, err := http.DefaultClient.Do(tokenCheckRequest)
	if err != nil {
		return nil, fmt.Errorf("token validation tokenCheckRequest error: %w", err)
	}

	return response, nil
}
