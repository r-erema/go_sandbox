package example4_test

import (
	"context"
	"github.com/Nerzal/gocloak/v11"
	"github.com/r-erema/go_sendbox/utils/test"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/tools/clientcmd"
	"testing"
)

// curl  -k -d "grant_type=password" -d "scope=openid" -d "client_id=account" -d "client_secret=dItyeR5mujp2O5qEkQ5eIMZwL44PzgoX" -d "username=admin" -d "password=admin"  https://localhost:8443/realms/master/protocol/openid-connect/token

const keyCloakURLEnvVar = "KEYCLOAK_URL"

func TestOIDC(t *testing.T) {
	/*cfg := oauth2.Config{
		ClientID:     "account",
		ClientSecret: "dItyeR5mujp2O5qEkQ5eIMZwL44PzgoX",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "",
			TokenURL: "https://localhost:8443/realms/master/protocol/openid-connect/token",
		},
	}
	token, err := cfg.PasswordCredentialsToken(context.Background(), "admin", "admin")
	require.NoError(t, err)
	_ = token*/
	// keycloakURL, ok := os.LookupEnv(keyCloakURLEnvVar)
	keycloakURL := "https://localhost:8443"
	//require.True(t, ok)

	client := gocloak.NewClient(keycloakURL, gocloak.SetAuthAdminRealms("admin/realms"), gocloak.SetAuthRealms("realms"))
	token, err := client.LoginAdmin(context.Background(), "admin", "admin", "master")
	require.NoError(t, err)

	cs, err := client.GetClients(context.Background(), token.AccessToken, "master", gocloak.GetClientsParams{})
	_ = cs
	c, err := client.GetClient(context.Background(), token.AccessToken, "master", "75c30f69-967f-4f5a-b3c6-0831966fbdfb")
	require.NoError(t, err)
	_ = c

	j, err := client.LoginClientTokenExchange(context.Background(), *c.ClientID, token.AccessToken, *c.Secret, "master", "", "")
	_ = j
	// client.GetClient()

	config, err := clientcmd.BuildConfigFromFlags("", *test.KubeConfigPtr(t))
	require.NoError(t, err)

	clientset, err := kubernetes.NewForConfig(config)
	require.NoError(t, err)

	pods, err := clientset.CoreV1().Pods("").List(context.Background(), v1.ListOptions{})
	require.NoError(t, err)

	_ = pods.Items[0].Status.PodIP
}
