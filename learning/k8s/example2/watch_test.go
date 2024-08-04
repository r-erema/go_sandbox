package example2_test

import (
	"context"
	"fmt"
	"mime"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/r-erema/go_sendbox/utils/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/streaming"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	versioned "k8s.io/client-go/rest/watch"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubectl/pkg/cmd"
	"k8s.io/kubectl/pkg/cmd/plugin"
)

func TestInformer(t *testing.T) {
	t.Parallel()

	cfg, err := clientcmd.BuildConfigFromFlags("", *test.KubeConfigPtr(t))
	require.NoError(t, err)

	kubeClientset, err := kubernetes.NewForConfig(cfg)
	require.NoError(t, err)

	opts := metaV1.ListOptions{
		TypeMeta:             metaV1.TypeMeta{Kind: "", APIVersion: ""},
		LabelSelector:        "",
		FieldSelector:        "",
		Watch:                true,
		AllowWatchBookmarks:  false,
		ResourceVersion:      "",
		ResourceVersionMatch: "",
		TimeoutSeconds:       nil,
		Limit:                0,
		Continue:             "",
	}

	cfg.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	cfg.GroupVersion = &v1.SchemeGroupVersion

	request := kubeClientset.CoreV1().RESTClient().Get().
		Resource("pods").
		VersionedParams(&opts, metaV1.ParameterCodec).
		Timeout(time.Second * 10)

	url := request.URL().String()

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
	require.NoError(t, err)

	RESTClient, ok := kubeClientset.CoreV1().RESTClient().(*rest.RESTClient)
	require.True(t, ok)

	resp, err := RESTClient.Client.Do(req)
	require.NoError(t, err)

	defer func() {
		err = resp.Body.Close()
		require.NoError(t, err)
	}()

	contentType := resp.Header.Get("Content-Type")
	mediaType, params, err := mime.ParseMediaType(contentType)
	require.NoError(t, err)

	negotiator := runtime.NewClientNegotiator(cfg.NegotiatedSerializer, *cfg.GroupVersion)
	objectDecoder, streamingSerializer, framer, err := negotiator.StreamDecoder(mediaType, params)
	require.NoError(t, err)

	framerReader := framer.NewFrameReader(resp.Body)

	watchEventDecoder := streaming.NewDecoder(framerReader, streamingSerializer)

	streamer := watch.NewStreamWatcher(
		versioned.NewDecoder(watchEventDecoder, objectDecoder),
		errors.NewClientErrorReporter(http.StatusInternalServerError, http.MethodGet, "ClientWatchDecoding"),
	)

	defaultConfigFlags := test.CLIConfigFlags(t)

	var mutex sync.Mutex

	go func(mutex *sync.Mutex) {
		time.Sleep(time.Millisecond * 100)

		mutex.Lock()
		defer mutex.Unlock()

		err = addTestPOD(defaultConfigFlags)
		assert.NoError(t, err)
	}(&mutex)

	event := <-streamer.ResultChan()

	assert.Equal(t, watch.Added, event.Type)
}

func addTestPOD(defaultConfigFlags *genericclioptions.ConfigFlags) error {
	command := cmd.NewDefaultKubectlCommandWithArgs(cmd.KubectlOptions{
		PluginHandler: cmd.NewDefaultPluginHandler(plugin.ValidPluginFilenamePrefixes),
		Arguments:     nil,
		ConfigFlags:   defaultConfigFlags,
		IOStreams:     genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr},
	})
	command.SetArgs([]string{"run", "test-busybox", "--image=busybox:1.35.0"})

	if err := command.Execute(); err != nil {
		return fmt.Errorf("command execution error: %w", err)
	}

	return nil
}
