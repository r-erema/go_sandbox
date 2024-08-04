package test

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/require"
)

const (
	etcdImage          = "quay.io/coreos/etcd:v3.5.1"
	kubeAPIServerImage = "k8s.gcr.io/kube-apiserver:v1.23.4"
	keycloakImage      = "quay.io/keycloak/keycloak:19.0.2"

	dockerStuffPath = "../../../docker/k8s"
	certPath        = "../../../docker/k8s/common_cert_for_all.crt"
	certKeyPath     = "../../../docker/k8s/common_cert_key_for_all.key"
)

func dockerClient(t *testing.T) *client.Client {
	t.Helper()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	require.NoError(t, err)

	return cli
}

func PullEtcdImage(t *testing.T) {
	t.Helper()

	cli := dockerClient(t)

	reader, err := cli.ImagePull(context.Background(), etcdImage, image.PullOptions{})
	defer func() {
		err = reader.Close()
		require.NoError(t, err)
	}()
	require.NoError(t, err)

	_, err = io.Copy(os.Stdout, reader)
	require.NoError(t, err)
}

func RunEtcdContainer(t *testing.T, hostPortBinding nat.PortBinding) string {
	t.Helper()

	cli := dockerClient(t)

	resp, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image: etcdImage,
		Cmd:   []string{"/usr/local/bin/etcd", "--advertise-client-urls=http://localhost:2379", "--listen-client-urls=http://0.0.0.0:2379"},
		ExposedPorts: nat.PortSet{
			"2379/tcp": struct{}{},
		},
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			"2379/tcp": []nat.PortBinding{hostPortBinding},
		},
	}, nil, nil, "")
	require.NoError(t, err)

	err = cli.ContainerStart(context.Background(), resp.ID, container.StartOptions{})
	require.NoError(t, err)

	return resp.ID
}

func PullKubeAPIServerImage(t *testing.T) {
	t.Helper()

	cli := dockerClient(t)

	reader, err := cli.ImagePull(context.Background(), kubeAPIServerImage, image.PullOptions{})
	defer func() {
		err = reader.Close()
		require.NoError(t, err)
	}()
	require.NoError(t, err)

	_, err = io.Copy(os.Stdout, reader)
	require.NoError(t, err)
}

func RunKubeAPIServer(t *testing.T, port, etcdHost, oidcIssuerURL string) string {
	t.Helper()

	cli := dockerClient(t)

	absDockerStuffPath, err := filepath.Abs(dockerStuffPath)
	require.NoError(t, err)

	resp, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image: kubeAPIServerImage,
		Cmd: []string{
			"kube-apiserver",
			"--secure-port=" + port,
			"--etcd-servers=" + etcdHost,
			"--service-account-signing-key-file=/var/run/kubernetes/common_cert_key_for_all.key",
			"--service-account-key-file=/var/run/kubernetes/common_cert_for_all.crt",
			"--service-account-issuer=https://kube.local",
			"--client-ca-file=/var/run/kubernetes/admin-auth.crt",
			"--tls-private-key-file=/var/run/kubernetes/common_cert_key_for_all.key",
			"--tls-cert-file=/var/run/kubernetes/common_cert_for_all.crt",
			"--audit-log-path=/var/run/kubernetes/audit.log",
			"--audit-policy-file=/var/run/kubernetes/audit-policy.yaml",
			"--oidc-issuer-url=" + oidcIssuerURL,
			"--oidc-client-id=k8s-auth-service",
			"--oidc-username-claim=name",
			"--oidc-groups-claim=groups",
			"--oidc-ca-file=/var/run/kubernetes/rootCA.crt",
			"--authorization-mode=RBAC",
		},
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: absDockerStuffPath,
				Target: "/var/run/kubernetes",
			},
		},
		ExtraHosts:  []string{"host.docker.internal:host-gateway"},
		NetworkMode: "host",
	}, nil, nil, "")
	require.NoError(t, err)

	err = cli.ContainerStart(context.Background(), resp.ID, container.StartOptions{})
	require.NoError(t, err)

	return resp.ID
}

func PullKeycloakImage(t *testing.T) {
	t.Helper()

	cli := dockerClient(t)

	reader, err := cli.ImagePull(context.Background(), keycloakImage, image.PullOptions{})
	defer func() {
		err = reader.Close()
		require.NoError(t, err)
	}()
	require.NoError(t, err)

	_, err = io.Copy(os.Stdout, reader)
	require.NoError(t, err)
}

func RunKeycloakContainer(t *testing.T, hostPortBinding nat.PortBinding) string {
	t.Helper()

	cli := dockerClient(t)

	absCertPath, err := filepath.Abs(certPath)
	require.NoError(t, err)
	absKeyCertPath, err := filepath.Abs(certKeyPath)
	require.NoError(t, err)

	resp, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image: keycloakImage,
		Env: []string{
			"KEYCLOAK_ADMIN=admin",
			"KEYCLOAK_ADMIN_PASSWORD=admin",
			"PROXY_ADDRESS_FORWARDING=true",
			"JAVA_OPTS=\"-Dkeycloak.profile.feature.token_exchange=enabled\"",
		},
		Cmd: []string{
			"start " +
				"--optimized " +
				"--hostname=localhost " +
				fmt.Sprintf("--hostname-port=%s ", hostPortBinding.HostPort) +
				"--https-certificate-file=/etc/x509/https/tls.crt " +
				"--https-certificate-key-file=/etc/x509/https/tls.key",
		},
		ExposedPorts: nat.PortSet{
			"8443/tcp": struct{}{},
		},
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: absCertPath,
				Target: "/etc/x509/https/tls.crt",
			},
			{
				Type:   mount.TypeBind,
				Source: absKeyCertPath,
				Target: "/etc/x509/https/tls.key",
			},
		},
		PortBindings: nat.PortMap{
			"8443/tcp": []nat.PortBinding{hostPortBinding},
		},
	}, nil, nil, "")
	require.NoError(t, err)

	err = cli.ContainerStart(context.Background(), resp.ID, container.StartOptions{})
	require.NoError(t, err)

	return resp.ID
}

func StopAndRemoveContainer(t *testing.T, containerID string) {
	t.Helper()

	cli := dockerClient(t)

	err := cli.ContainerStop(context.Background(), containerID, container.StopOptions{})
	require.NoError(t, err)

	err = cli.ContainerRemove(context.Background(), containerID, container.RemoveOptions{
		RemoveVolumes: true,
	})
	require.NoError(t, err)
}
