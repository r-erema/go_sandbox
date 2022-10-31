package example5_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/Nerzal/gocloak/v11"
	"github.com/docker/go-connections/nat"
	"github.com/phayes/freeport"
	"github.com/r-erema/go_sendbox/utils/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubectl/pkg/cmd/util"
)

const (
	roleManifestPath = "./role.yaml"

	keycloakRealm = "k8s-realm"

	keycloakAdminUsername = "admin"
	keycloakAdminPassword = "admin"

	keycloakClientID        = "k8s-auth-service"
	keycloakIDOfClient      = "e76ec143-f2f8-48c6-ba01-5ac51113be85"
	keycloakClientScopeName = "groups"

	keycloakUserUsername = "k8s-admin"
	keycloakUserPassword = "123"
	keycloakUserName     = "John Doe"
	keycloakUserGroup    = "developers"

	k8sClusterName          = "OIDCTestCluster"
	k8sAdminUser            = "Admin"
	k8sAdminContext         = "admin-context"
	k8sCAAuthorityPath      = "../../../docker/k8s/rootCA.crt"
	k8sAdminAuthCertPath    = "../../../docker/k8s/admin-auth.crt"
	k8sAdminAuthCertKeyPath = "../../../docker/k8s/admin-auth.key"
	k8sTestUser             = "JohnDoe"
	k8sTestContext          = "developer-context"
)

func TestOIDC(t *testing.T) { //nolint: paralleltest, tparallel
	tests := []struct {
		name        string
		expectError bool
	}{
		{
			name:        "Expect error",
			expectError: true,
		},
		{
			name:        "Error isn't expected",
			expectError: false,
		},
	}

	var mutex sync.Mutex

	mutex.Lock()
	pullImages(t)
	mutex.Unlock()

	for _, tt := range tests {
		testCase := tt

		util.BehaviorOnFatal(func(msg string, _ int) {
			if testCase.expectError {
				require.Equal(t, "error: You must be logged in to the server (Unauthorized)\n", msg)
			}
		})

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			mutex.Lock()
			KeycloakURL, IDPIssuerURL, KubeConfigPath, clean := prepareInfrastructure(t)
			defer clean()

			stepPrepareClusterRole(t, string(KubeConfigPath))
			time.Sleep(time.Second * 10)

			keycloakUser := defaultKeycloakUser()
			if !testCase.expectError {
				keycloakUser.Groups = &[]string{keycloakUserGroup}
			}

			time.Sleep(time.Second * 5)
			clientID, secret, refreshToken, IDToken := stepPrepareKeycloakClient(
				t,
				string(KeycloakURL),
				keycloakRealmRepresentation(
					ssoSessionIdleTimeout(gocloak.IntP(1)),
					users(&[]gocloak.User{keycloakUser}),
				),
			)

			stepPrepareKubectlDeveloperContext(
				t,
				string(KubeConfigPath),
				string(IDPIssuerURL),
				clientID,
				secret,
				refreshToken,
				IDToken,
			)

			err := test.RunKubectlCommand(test.DefaultConfigFlags(), []string{
				"--context=" + k8sTestContext,
				"get",
				"pods",
				"--kubeconfig=" + string(KubeConfigPath),
			})
			require.NoError(t, err)
			mutex.Unlock()
		})
	}
}

func TestNotExpiredRefreshToken(t *testing.T) { //nolint: paralleltest
	KeycloakURL, IDPIssuerURL, KubeConfigPath, clean := prepareInfrastructure(t)
	defer clean()

	keycloakUser := defaultKeycloakUser()
	keycloakUser.Groups = &[]string{keycloakUserGroup}

	const (
		accessTokenLifetimeSeconds  = 1
		refreshTokenLifetimeSeconds = 60
	)

	time.Sleep(time.Second * 11)

	clientID, secret, refreshToken, IDToken := stepPrepareKeycloakClient(
		t,
		string(KeycloakURL),
		keycloakRealmRepresentation(
			accessTokenLifespan(gocloak.IntP(accessTokenLifetimeSeconds)),
			ssoSessionIdleTimeout(gocloak.IntP(refreshTokenLifetimeSeconds)),
			users(&[]gocloak.User{keycloakUser}),
		),
	)

	stepPrepareClusterRole(t, string(KubeConfigPath))
	time.Sleep(time.Second * 10)

	stepPrepareKubectlDeveloperContext(
		t,
		string(KubeConfigPath),
		string(IDPIssuerURL),
		clientID,
		secret,
		refreshToken,
		IDToken,
	)

	cfg, err := clientcmd.BuildConfigFromFlags("", string(KubeConfigPath))
	require.NoError(t, err)

	kubeClientset, err := kubernetes.NewForConfig(cfg)
	require.NoError(t, err)

	list, err := kubeClientset.CoreV1().Pods("").List(context.Background(), v1.ListOptions{})
	require.NoError(t, err)
	assert.Empty(t, list.Items)
}

func TestProveBug_UpdatedGoodRefreshTokenInsteadOfExpiredOneIsNotBeingApplied(t *testing.T) { //nolint: paralleltest
	KeycloakURL, IDPIssuerURL, KubeConfigPath, clean := prepareInfrastructure(t)
	defer clean()

	keycloakUser := defaultKeycloakUser()
	keycloakUser.Groups = &[]string{keycloakUserGroup}

	const (
		accessTokenLifetimeSeconds       = 1
		refreshTokenShortLifetimeSeconds = 2
		refreshTokenLongLifetimeSeconds  = 30
	)

	time.Sleep(time.Second * 11)

	clientID, secret, refreshToken, IDToken := stepPrepareKeycloakClient(
		t,
		string(KeycloakURL),
		keycloakRealmRepresentation(
			accessTokenLifespan(gocloak.IntP(accessTokenLifetimeSeconds)),
			ssoSessionIdleTimeout(gocloak.IntP(refreshTokenShortLifetimeSeconds)),
			users(&[]gocloak.User{keycloakUser}),
		),
	)

	stepPrepareClusterRole(t, string(KubeConfigPath))
	time.Sleep(time.Second * 10)

	stepPrepareKubectlDeveloperContext(
		t,
		string(KubeConfigPath),
		string(IDPIssuerURL),
		clientID,
		secret,
		refreshToken,
		IDToken,
	)

	cfg, err := clientcmd.BuildConfigFromFlags("", string(KubeConfigPath))
	require.NoError(t, err)

	kubeClientset, err := kubernetes.NewForConfig(cfg)
	require.NoError(t, err)

	stepIntentionallyExpireTokensAndFailRequest := func() {
		stepSleepToExpireToken(accessTokenLifetimeSeconds)
		stepSleepToExpireToken(refreshTokenShortLifetimeSeconds)

		list, err := kubeClientset.CoreV1().Pods("").List(context.Background(), v1.ListOptions{})
		assert.Error(t, err)
		assert.Nil(t, list.Items)
	}

	stepIntentionallyExpireTokensAndFailRequest()

	stepAssertNotAuthorizedEvenWithGoodRefreshToken := func() {
		keycloakAdmin := gocloak.NewClient(string(KeycloakURL), gocloak.SetAuthAdminRealms("admin/realms"), gocloak.SetAuthRealms("realms"))
		adminToken, err := keycloakAdmin.LoginAdmin(context.Background(), keycloakAdminUsername, keycloakAdminPassword, "master")
		require.NoError(t, err)
		err = keycloakAdmin.UpdateRealm(context.Background(), adminToken.AccessToken, *keycloakRealmRepresentation(
			ssoSessionIdleTimeout(gocloak.IntP(refreshTokenLongLifetimeSeconds)),
		))
		require.NoError(t, err)
		clientSecret, err := keycloakAdmin.GetClientSecret(context.Background(), adminToken.AccessToken, keycloakRealm, keycloakIDOfClient)
		require.NoError(t, err)
		token, err := keycloakAdmin.GetToken(context.Background(), keycloakRealm, gocloak.TokenOptions{
			ClientID:     gocloak.StringP(keycloakClientID),
			ClientSecret: clientSecret.Value,
			GrantType:    gocloak.StringP("password"),
			Scope:        gocloak.StringP("openid"),
			Username:     gocloak.StringP(keycloakUserUsername),
			Password:     gocloak.StringP(keycloakUserPassword),
		})
		require.NoError(t, err)
		stepPrepareKubectlDeveloperContext(
			t,
			string(KubeConfigPath),
			string(IDPIssuerURL),
			clientID,
			secret,
			token.RefreshToken,
			token.IDToken,
		)

		list, err := kubeClientset.CoreV1().Pods("").List(context.Background(), v1.ListOptions{})
		assert.Error(t, err)
		assert.Nil(t, list.Items)
	}
	stepAssertNotAuthorizedEvenWithGoodRefreshToken()
}

func stepPrepareKeycloakClient( //nolint:nonamedreturns
	t *testing.T,
	keycloakURL string,
	realmRepresentation *gocloak.RealmRepresentation,
) (clientID, secret, refreshToken, idToken string) {
	t.Helper()

	keycloakAdmin := gocloak.NewClient(keycloakURL, gocloak.SetAuthAdminRealms("admin/realms"), gocloak.SetAuthRealms("realms"))
	adminToken, err := keycloakAdmin.LoginAdmin(context.Background(), keycloakAdminUsername, keycloakAdminPassword, "master")
	require.NoError(t, err)

	realm, err := keycloakAdmin.CreateRealm(context.Background(), adminToken.AccessToken, *realmRepresentation)
	require.NoError(t, err)
	require.Equal(t, realm, keycloakRealm)

	clientSecret, err := keycloakAdmin.GetClientSecret(context.Background(), adminToken.AccessToken, keycloakRealm, keycloakIDOfClient)
	require.NoError(t, err)

	token, err := keycloakAdmin.GetToken(context.Background(), keycloakRealm, gocloak.TokenOptions{
		ClientID:     gocloak.StringP(keycloakClientID),
		ClientSecret: clientSecret.Value,
		GrantType:    gocloak.StringP("password"),
		Scope:        gocloak.StringP("openid"),
		Username:     gocloak.StringP(keycloakUserUsername),
		Password:     gocloak.StringP(keycloakUserPassword),
	})
	require.NoError(t, err)

	return keycloakClientID, *clientSecret.Value, token.RefreshToken, token.IDToken
}

func stepPrepareKubectlDeveloperContext(
	t *testing.T,
	kubeConfigPath,
	idpIssuerURL,
	clientID,
	secret,
	refreshToken,
	idToken string,
) {
	t.Helper()

	err := test.RunKubectlCommand(test.DefaultConfigFlags(), []string{
		"config",
		"set-credentials",
		k8sTestUser,
		"--auth-provider=oidc",
		"--auth-provider-arg=idp-issuer-url=" + idpIssuerURL,
		"--auth-provider-arg=client-id=" + clientID,
		"--auth-provider-arg=client-secret=" + secret,
		"--auth-provider-arg=refresh-token=" + refreshToken,
		"--auth-provider-arg=id-token=" + idToken,
		"--kubeconfig=" + kubeConfigPath,
	})
	require.NoError(t, err)

	err = test.RunKubectlCommand(test.DefaultConfigFlags(), []string{
		"config",
		"set-context",
		k8sTestContext,
		"--cluster=" + k8sClusterName,
		"--user=" + k8sTestUser,
		"--kubeconfig=" + kubeConfigPath,
	})
	require.NoError(t, err)

	err = test.RunKubectlCommand(test.DefaultConfigFlags(), []string{
		"config",
		"use-context",
		k8sTestContext,
		"--kubeconfig=" + kubeConfigPath,
	})
	require.NoError(t, err)
}

func stepPrepareClusterRole(t *testing.T, kubeConfigPath string) {
	t.Helper()

	err := test.RunKubectlCommand(test.DefaultConfigFlags(), []string{
		"apply",
		"-f",
		roleManifestPath,
		"--kubeconfig=" + kubeConfigPath,
	})
	require.NoError(t, err)
}

func stepSleepToExpireToken(duration time.Duration) {
	time.Sleep(duration)
}

func runDockerContainers(t *testing.T, etcdHostPort, keycloakHostPort, kubeAPIServerPort string) []string {
	t.Helper()

	containerIDs := make([]string, 3)

	etcdContainerID := test.RunEtcdContainer(t, nat.PortBinding{HostIP: "0.0.0.0", HostPort: etcdHostPort})
	containerIDs[0] = etcdContainerID

	keycloakContainerID := test.RunKeycloakContainer(t, nat.PortBinding{HostIP: "0.0.0.0", HostPort: keycloakHostPort})

	func() {
		ticker, attempt := time.NewTicker(time.Second), 0
		defer func() {
			ticker.Stop()
		}()

		var resp *http.Response
		defer func() {
			if resp != nil {
				err := resp.Body.Close()
				require.NoError(t, err)
			}
		}()

		for range ticker.C {
			attempt++

			if attempt > 10 {
				return
			}

			req, err := http.NewRequestWithContext(
				context.Background(),
				http.MethodGet,
				fmt.Sprintf("https://localhost:%s/realms/master", keycloakHostPort),
				http.NoBody,
			)
			require.NoError(t, err)

			resp, err = http.DefaultClient.Do(req)
			if err != nil {
				continue
			}

			if resp.StatusCode == http.StatusOK {
				return
			}
		}
	}()

	containerIDs[1] = keycloakContainerID

	kubeAPIServerContainerID := test.RunKubeAPIServer(
		t,
		kubeAPIServerPort,
		fmt.Sprintf("host.docker.internal:%s", etcdHostPort),
		fmt.Sprintf("https://localhost:%s/realms/k8s-realm", keycloakHostPort),
	)
	containerIDs[2] = kubeAPIServerContainerID

	return containerIDs
}

func keycloakRealmRepresentation(realmOptions ...func(*gocloak.RealmRepresentation)) *gocloak.RealmRepresentation {
	realm := &gocloak.RealmRepresentation{
		Realm:   gocloak.StringP(keycloakRealm),
		Enabled: gocloak.BoolP(true),
		Clients: &[]gocloak.Client{
			{
				ClientID:                     gocloak.StringP(keycloakClientID),
				ID:                           gocloak.StringP(keycloakIDOfClient),
				PublicClient:                 gocloak.BoolP(false),
				AuthorizationServicesEnabled: gocloak.BoolP(false),
				ServiceAccountsEnabled:       gocloak.BoolP(false),
				DirectAccessGrantsEnabled:    gocloak.BoolP(true),
				DefaultClientScopes:          &[]string{keycloakClientScopeName},
			},
		},
		ClientScopes: &[]gocloak.ClientScope{
			{
				Name:     gocloak.StringP(keycloakClientScopeName),
				Protocol: gocloak.StringP("openid-connect"),
				ClientScopeAttributes: &gocloak.ClientScopeAttributes{
					DisplayOnConsentScreen: gocloak.StringP("true"),
				},
				ProtocolMappers: &[]gocloak.ProtocolMappers{
					{
						Name:           gocloak.StringP("name"),
						Protocol:       gocloak.StringP("openid-connect"),
						ProtocolMapper: gocloak.StringP("oidc-usermodel-attribute-mapper"),
						ProtocolMappersConfig: &gocloak.ProtocolMappersConfig{
							AccessTokenClaim:   gocloak.StringP("true"),
							ClaimName:          gocloak.StringP("name"),
							IDTokenClaim:       gocloak.StringP("true"),
							JSONTypeLabel:      gocloak.StringP("string"),
							UserAttribute:      gocloak.StringP("name"),
							UserinfoTokenClaim: gocloak.StringP("name"),
						},
					},
					{
						Name:           gocloak.StringP("groups"),
						Protocol:       gocloak.StringP("openid-connect"),
						ProtocolMapper: gocloak.StringP("oidc-group-membership-mapper"),
						ProtocolMappersConfig: &gocloak.ProtocolMappersConfig{
							AccessTokenClaim:   gocloak.StringP("true"),
							ClaimName:          gocloak.StringP("groups"),
							FullPath:           gocloak.StringP("false"),
							IDTokenClaim:       gocloak.StringP("true"),
							UserinfoTokenClaim: gocloak.StringP("true"),
						},
					},
				},
			},
		},
		Groups: &[]interface{}{
			gocloak.Group{Name: gocloak.StringP(keycloakUserGroup)},
		},
	}

	for _, option := range realmOptions {
		option(realm)
	}

	return realm
}

func accessTokenLifespan(timeout *int) func(*gocloak.RealmRepresentation) {
	return func(r *gocloak.RealmRepresentation) {
		r.AccessTokenLifespan = timeout
	}
}

func ssoSessionIdleTimeout(timeout *int) func(*gocloak.RealmRepresentation) {
	return func(r *gocloak.RealmRepresentation) {
		r.SsoSessionIdleTimeout = timeout
	}
}

func users(users *[]gocloak.User) func(*gocloak.RealmRepresentation) {
	return func(r *gocloak.RealmRepresentation) {
		r.Users = users
	}
}

func defaultKeycloakUser() gocloak.User {
	return gocloak.User{
		Username: gocloak.StringP(keycloakUserUsername),
		Credentials: &[]gocloak.CredentialRepresentation{
			{
				Type:      gocloak.StringP("password"),
				Value:     gocloak.StringP(keycloakUserPassword),
				Temporary: gocloak.BoolP(false),
			},
		},
		Enabled:    gocloak.BoolP(true),
		Attributes: &map[string][]string{"name": {keycloakUserName}},
	}
}

func pullImages(t *testing.T) {
	t.Helper()

	var waitGroup sync.WaitGroup

	waitGroup.Add(3)

	go func() {
		test.PullEtcdImage(t)
		waitGroup.Done()
	}()

	go func() {
		test.PullKeycloakImage(t)
		waitGroup.Done()
	}()

	go func() {
		test.PullKubeAPIServerImage(t)
		waitGroup.Done()
	}()

	waitGroup.Wait()
}

type (
	keycloakURL                  string
	kubeConfigPath               string
	idpIssuerURL                 string
	removeInfrastructureCallback func()
)

func prepareInfrastructure(t *testing.T) (keycloakURL, idpIssuerURL, kubeConfigPath, removeInfrastructureCallback) {
	t.Helper()

	port, err := freeport.GetFreePort()
	require.NoError(t, err)

	etcdHostPort := strconv.Itoa(port)
	port, err = freeport.GetFreePort()
	require.NoError(t, err)

	keycloakHostPort := strconv.Itoa(port)
	port, err = freeport.GetFreePort()
	require.NoError(t, err)

	kubeAPIServerPort := strconv.Itoa(port)

	containerIDs := runDockerContainers(t, etcdHostPort, keycloakHostPort, kubeAPIServerPort)

	KubeConfigPath, err := filepath.Abs(fmt.Sprintf("./kube_cfg_%d", time.Now().UnixNano()))
	require.NoError(t, err)

	test.PrepareKubeConfigContext(t,
		KubeConfigPath,
		k8sClusterName,
		k8sAdminUser,
		k8sAdminContext,
		"https://localhost:"+kubeAPIServerPort,
		k8sCAAuthorityPath,
		k8sAdminAuthCertPath,
		k8sAdminAuthCertKeyPath,
	)

	KeycloakURL := fmt.Sprintf("https://localhost:%s", keycloakHostPort)
	IDPIssuerURL := fmt.Sprintf("%s/realms/%s", KeycloakURL, keycloakRealm)

	return keycloakURL(KeycloakURL),
		idpIssuerURL(IDPIssuerURL),
		kubeConfigPath(KubeConfigPath),
		func() {
			for _, id := range containerIDs {
				test.StopAndRemoveContainer(t, id)
			}

			err = os.Remove(KubeConfigPath)
			require.NoError(t, err)
		}
}
