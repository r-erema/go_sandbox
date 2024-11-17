package prometheus_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/phayes/freeport"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"
	"github.com/prometheus/prometheus/discovery"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"github.com/r-erema/go_sendbox/utils/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

const testMetricName = "test_metric"

func TestGauge(t *testing.T) {
	t.Parallel()

	exporterPort, err := freeport.GetFreePort()
	require.NoError(t, err)

	go func() {
		err = RunExporter(exporterPort)
		assert.NoError(t, err)
	}()
	t.Logf("Exporter started on port: http://localhost:%d", exporterPort)

	cfg := preparePrometheusConfig(t, model.LabelValue(fmt.Sprintf("host.docker.internal:%d", exporterPort)))
	defer func() {
		err = cfg.Close()
		require.NoError(t, err)
	}()

	prometheusPort, err := freeport.GetFreePort()
	require.NoError(t, err)

	prometheusURL := fmt.Sprintf("http://localhost:%d", prometheusPort)
	t.Logf("Prometheus started on port: %s", prometheusURL)
	containerID := test.RunPrometheusContainer(t, nat.PortBinding{HostIP: "0.0.0.0", HostPort: strconv.Itoa(prometheusPort)}, cfg.Name())
	t.Cleanup(func() {
		test.StopAndRemoveContainer(t, containerID)
	})

	client, err := api.NewClient(api.Config{Address: prometheusURL})
	require.NoError(t, err)
	waitPrometheusReady(t, client)
	promAPI := v1.NewAPI(client)

	value, _, err := promAPI.Query(context.Background(), fmt.Sprintf("{__name__=%q}", testMetricName), time.Now())
	require.NoError(t, err)

	maxAttempts := 10
	for value.String() == "" && maxAttempts > 0 {
		time.Sleep(time.Second)

		value, _, err = promAPI.Query(context.Background(), fmt.Sprintf("{__name__=%q}", testMetricName), time.Now())
		require.NoError(t, err)

		maxAttempts--
	}

	assert.Contains(t, value.String(), testMetricName)
}

func preparePrometheusConfig(t *testing.T, scraperPath model.LabelValue) *os.File {
	t.Helper()

	cfgFile, err := os.OpenFile(filepath.Join(t.TempDir(), "prometheus.yml"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
	require.NoError(t, err)

	cfg := config.Config{ScrapeConfigs: []*config.ScrapeConfig{
		{
			JobName: "test_exporter",
			ServiceDiscoveryConfigs: discovery.Configs{
				discovery.StaticConfig{
					&targetgroup.Group{
						Targets: []model.LabelSet{{"__address__": scraperPath}},
					},
				},
			},
			ScrapeInterval: model.Duration(time.Second),
		},
	}}
	enc := yaml.NewEncoder(cfgFile)
	err = enc.Encode(cfg)
	require.NoError(t, err)

	return cfgFile
}

func RunExporter(port int) error {
	vec := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: testMetricName,
		Help: "test metric",
	}, []string{"test_label_key"})

	if err := prometheus.Register(vec); err != nil {
		return fmt.Errorf("could not register test metric: %w", err)
	}

	metric, err := vec.GetMetricWith(prometheus.Labels{"test_label_key": "test_label_value"})
	if err != nil {
		return fmt.Errorf("could not get test metric: %w", err)
	}

	tick := 1

	go func() {
		for range time.NewTicker(time.Second).C {
			log.Printf("Tick: %d", tick)
			metric.Add(float64(tick))

			tick++
		}
	}()

	http.Handle("/metrics", promhttp.Handler())

	log.Printf("Start listening: %d", port)

	if err = (&http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           http.DefaultServeMux,
		ReadHeaderTimeout: time.Second,
	}).ListenAndServe(); err != nil {
		return fmt.Errorf("could not start exporter: %w", err)
	}

	return nil
}

func waitPrometheusReady(t *testing.T, client api.Client) {
	t.Helper()

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		client.URL("/-/ready", make(map[string]string)).String(),
		http.NoBody,
	)
	require.NoError(t, err)

	callWaitTime := time.Millisecond * 100
	for range time.Tick(callWaitTime) {
		resp, _, err := client.Do(context.Background(), req)
		if err != nil {
			if errors.Is(err, syscall.ECONNRESET) {
				continue
			}

			require.NoError(t, err)
		}

		err = resp.Body.Close()
		require.NoError(t, err)

		if resp.StatusCode == http.StatusOK {
			return
		}
	}
}
