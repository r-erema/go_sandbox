package example1

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-oauth2/oauth2/v4/server"
	"k8s.io/apimachinery/pkg/util/json"
)

type AuthorizationServer struct {
	oAuthHandler *server.Server
}

func NewAuthorizationServer(oAuthHandler *server.Server) *AuthorizationServer {
	return &AuthorizationServer{oAuthHandler: oAuthHandler}
}

func (at AuthorizationServer) Run(addr string) error {
	router := http.NewServeMux()

	router.Handle("/token", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if err := at.oAuthHandler.HandleTokenRequest(writer, request); err != nil {
			log.Printf("handling token request error: %s", err)

			http.Error(writer, "", http.StatusUnauthorized)
		}
	}))

	router.Handle("/interception-endpoint", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		tokenInfo, err := at.oAuthHandler.ValidationBearerToken(request)
		if err != nil {
			log.Printf("validation token request error: %s", err)
			http.Error(writer, "", http.StatusUnauthorized)

			return
		}

		body, err := json.Marshal(tokenInfo)
		if err != nil {
			log.Printf("marshaling token info error: %s", err)
			http.Error(writer, "", http.StatusInternalServerError)

			return
		}

		_, err = writer.Write(body)
		if err != nil {
			log.Printf("writing body error: %s", err)
			http.Error(writer, "", http.StatusInternalServerError)

			return
		}
	}))

	if err := (&http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: time.Second,
	}).ListenAndServe(); err != nil {
		return fmt.Errorf("authorization server listen and serve error: %w", err)
	}

	return nil
}
