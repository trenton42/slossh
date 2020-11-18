package httprecorder

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/spf13/pflag"
	"github.com/trenton42/slossh/pkg/recorders"
	"github.com/trenton42/slossh/pkg/session"
)

func init() {
	hr := HTTPRecorder{}
	hr.client = http.DefaultClient

	recorders.Register(&hr)
}

// HTTPRecorder records ssh attempts into a file as json objects
type HTTPRecorder struct {
	client *http.Client
	url    string
}

// Name returns the name of this recorder
func (h *HTTPRecorder) Name() string {
	return "http"
}

// Options returns the flags that can configure this recorder
func (h *HTTPRecorder) Options() *pflag.FlagSet {
	flags := pflag.NewFlagSet("http", pflag.ContinueOnError)
	flags.StringVar(&h.url, "http-url", "", "URL to send post requests to")
	return flags
}

// Record records sessions
func (h *HTTPRecorder) Record(s session.SlosshSession) {
	if h.url == "" {
		return
	}
	data, err := json.Marshal(s)
	if err != nil {
		return
	}
	res, err := h.client.Post(h.url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return
	}
	res.Body.Close()
}

// Close closes this recorder and any associated resources
func (h *HTTPRecorder) Close() {
}
