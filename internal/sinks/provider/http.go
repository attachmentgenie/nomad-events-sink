package provider

import (
	"encoding/json"
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/grafana/dskit/backoff"
	"github.com/grafana/dskit/flagext"
	prom_api "github.com/grafana/loki/clients/pkg/promtail/api"
	"github.com/grafana/loki/clients/pkg/promtail/client"
	"github.com/grafana/loki/pkg/logproto"
	lokiflag "github.com/grafana/loki/pkg/util/flagext"
	nomad_api "github.com/hashicorp/nomad/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	"github.com/sirupsen/logrus"
)

// HTTPManager represents the various methods for interacting with Pigeon.
type HTTPManager struct {
	client  client.Client
	rootURL string
	log     *logrus.Logger
}

type HTTPOpts struct {
	ExternalLabels lokiflag.LabelSet
	Log            *logrus.Logger
	Password       string
	RootURL        string
	Timeout        time.Duration
	Username       string
}

// NewHTTP initializes a HTTP notification dispatcher object.
func NewHTTP(opts HTTPOpts) (*HTTPManager, error) {
	serverURL := flagext.URLValue{}
	err := serverURL.Set(opts.RootURL)
	if err != nil {
		return nil, err
	}
	cfg := client.Config{
		URL:                    serverURL,
		BatchWait:              100 * time.Millisecond,
		BatchSize:              10,
		DropRateLimitedBatches: true,
		Client:                 config.HTTPClientConfig{BasicAuth: &config.BasicAuth{Username: opts.Username, Password: config.Secret(opts.Password)}},
		BackoffConfig:          backoff.Config{MinBackoff: 1 * time.Millisecond, MaxBackoff: 2 * time.Millisecond, MaxRetries: 3},
		ExternalLabels:         opts.ExternalLabels,
		Timeout:                opts.Timeout,
		TenantID:               "tenant-default",
	}

	reg := prometheus.NewRegistry()
	m := client.NewMetrics(reg)
	c, err := client.New(m, cfg, 0, 999999, false, log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr)))
	if err != nil {
		//fmt.Errorf("sink HTTP provider misconfigured")
		return nil, err
	}

	httpMgr := &HTTPManager{
		client:  c,
		rootURL: opts.RootURL,
		log:     opts.Log,
	}

	return httpMgr, nil
}

// Push sends out events to an HTTP Endpoint.
func (m *HTTPManager) Push(events []nomad_api.Event) {
	for _, event := range events {
		labelSet := model.LabelSet{
			"topic": model.LabelValue(event.Topic),
			"type":  model.LabelValue(event.Type),
		}

		event_info, _ := json.Marshal(event)
		logEntry := prom_api.Entry{Labels: labelSet, Entry: logproto.Entry{Timestamp: time.Now(), Line: string(event_info)}}
		m.client.Chan() <- logEntry
	}
}

// Name returns the notification provider name.
func (m *HTTPManager) Name() string {
	return "http"
}
